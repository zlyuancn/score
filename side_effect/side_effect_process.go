package side_effect

import (
	"context"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/zly-app/zapp/component/gpool"
	"github.com/zly-app/zapp/logger"
	"go.uber.org/zap"

	"github.com/zlyuancn/score/dao"
	"github.com/zlyuancn/score/model"
)

var sideEffectTypeResolver = map[model.SideEffectType]func(ctx context.Context, data *model.SideEffect) error{
	model.SideEffectType_ScoreChange: scoreChangeHandle,
}

// 副作用补偿
func compensationSideEffect(ctx context.Context, payload string) error {
	data := &model.SideEffect{}
	err := sonic.UnmarshalString(payload, data)
	if err != nil {
		logger.Error(ctx, "compensationSideEffect call UnmarshalString data fail.", zap.Any("payload", payload), zap.Error(err))
		return nil
	}

	r, ok := sideEffectTypeResolver[data.Type]
	if ok {
		return r(ctx, data)
	}
	return fmt.Errorf("compensationSideEffect got not supported type=%d", data.Type)
}

// 准备副作用
func PrepareSideEffect(ctx context.Context, data *model.SideEffect) error {
	payload, err := sonic.MarshalString(data)
	if err != nil {
		logger.Error(ctx, "PrepareSideEffect call MarshalString data fail.", zap.Any("data", data), zap.Error(err))
		return err
	}

	err = mqTool.Send(ctx, payload)
	if err != nil {
		logger.Error(ctx, "PrepareSideEffect call mqTool.Send fail.", zap.String("payload", payload), zap.Error(err))
		return err
	}
	return nil
}

// 处理副作用
func processSideEffect(ctx context.Context, st *model.ScoreType, orderID, uid string,
	t model.SideEffectType, fn func(seName string, se SideEffect) error) error {
	seList := seMap[t]
	if len(seList) == 0 {
		return nil
	}

	fns := make([]func() error, 0, len(seList))
	for k, v := range seList {
		name, se := k, v
		fns = append(fns, func() error {
			// 获取订单副作用状态
			ok, err := dao.GetOrderSideEffectStatus(ctx, orderID, uid, name, int(t))
			if err != nil {
				logger.Error(ctx, "processSideEffect call GetOrderSideEffectStatus fail.", zap.Int("SideNameType", int(t)), zap.String("SideEffectName", name), zap.Error(err))
				return err
			}
			if ok {
				return nil
			}

			err = fn(name, se)
			if err != nil {
				logger.Error(ctx, "processSideEffect call fail.", zap.Int("SideNameType", int(t)), zap.String("SideEffectName", name), zap.Error(err))
				return err
			}

			// 标记订单副作用状态已完成
			err = dao.MarkOrderSideEffectStatusOk(ctx, orderID, uid, name, int(t), int64(st.OrderStatusExpireDay)*86400)
			if err != nil {
				logger.Error(ctx, "processSideEffect dao.MarkOrderFlowStatusOk fail.", zap.Int("SideNameType", int(t)), zap.String("SideEffectName", name), zap.Error(err))
				// 这里不影响主进程
			}

			return nil
		})
	}

	err := gpool.GetDefGPool().GoAndWait(fns...)
	if err != nil {
		logger.Error(ctx, "processSideEffect fail.", zap.Int("SideNameType", int(t)), zap.Error(err))
		return err
	}
	return nil
}
