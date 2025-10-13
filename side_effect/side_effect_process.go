package side_effect

import (
	"context"
	"fmt"
	"strconv"

	"github.com/bytedance/sonic"
	"github.com/zly-app/zapp/component/gpool"
	"github.com/zly-app/zapp/filter"
	"github.com/zly-app/zapp/logger"
	"github.com/zly-app/zapp/pkg/utils"
	"go.uber.org/zap"

	"github.com/zlyuancn/score/dao"
	"github.com/zlyuancn/score/model"
	"github.com/zlyuancn/score/score_type"
)

var sideEffectTypeResolver = map[model.SideEffectType]func(ctx context.Context, st *model.ScoreType, data *model.SideEffectData) error{
	model.SideEffectType_AfterScoreChange: afterScoreChangeHandle,
}

type SideEffectProcess func(ctx context.Context, seName string, se SideEffect, st *model.ScoreType, data *model.SideEffectData) error

// 副作用补偿
func compensationSideEffect(ctx context.Context, payload string) error {
	data := &model.SideEffectData{}
	err := sonic.UnmarshalString(payload, data)
	if err != nil {
		logger.Error(ctx, "compensationSideEffect call UnmarshalString data fail.", zap.Any("payload", payload), zap.Error(err))
		return nil
	}

	// 检查积分类型
	st, err := score_type.ForceGetScoreType(ctx, data.ScoreTypeID)
	if err != nil {
		return err
	}

	r, ok := sideEffectTypeResolver[data.Type]
	if ok {
		return r(ctx, st, data)
	}
	return fmt.Errorf("compensationSideEffect got not supported type=%d", data.Type)
}

// 为副作用添加一个守护程序, 当副作用处理失败后会自动重试
func AddSideEffectDaemon(ctx context.Context, data *model.SideEffectData) error {
	payload, err := sonic.MarshalString(data)
	if err != nil {
		logger.Error(ctx, "AddSideEffectDaemon call MarshalString data fail.", zap.Any("data", data), zap.Error(err))
		return err
	}

	err = mqTool.Send(ctx, payload)
	if err != nil {
		logger.Error(ctx, "AddSideEffectDaemon call mqTool.Send fail.", zap.String("payload", payload), zap.Error(err))
		return err
	}
	return nil
}

type triggerSideEffectAppFilterReq struct {
	Data *model.SideEffectData
	Name string
}

/*
立即触发副作用

	data 副作用数据
	t 副作用类型
	fn 如何处理副作用
*/
func TriggerSideEffect(ctx context.Context, data *model.SideEffectData, fn SideEffectProcess) error {

	// 检查积分类型
	st, err := score_type.GetScoreType(ctx, data.ScoreTypeID)
	if err != nil {
		logger.Error(ctx, "TriggerSideEffect call GetScoreType fail.", zap.Int("SideNameType", int(data.Type)), zap.Any("data", data), zap.Error(err))
		return err
	}

	seList := seMap[data.Type]
	if len(seList) == 0 {
		return nil
	}

	ctx = utils.Otel.CtxStart(ctx, "TriggerSideEffect")
	defer utils.Otel.CtxEnd(ctx)

	fns := make([]func() error, 0, len(seList))
	for k, v := range seList {
		name, se := k, v
		fns = append(fns, func() error {
			// 获取订单副作用状态
			ok, err := dao.GetOrderSideEffectStatus(ctx, data.OrderID, data.Uid, name, int(data.Type))
			if err != nil {
				logger.Error(ctx, "TriggerSideEffect call GetOrderSideEffectStatus fail.", zap.Int("SideNameType", int(data.Type)), zap.String("SideEffectName", name), zap.Any("data", data), zap.Error(err))
				return err
			}
			if ok {
				return nil
			}

			ctx, chain := filter.GetClientFilter(ctx, "TriggerSideEffect", strconv.FormatInt(int64(data.Type), 10), name)
			r := &triggerSideEffectAppFilterReq{
				Data: data,
				Name: name,
			}
			_, err = chain.Handle(ctx, r, func(ctx context.Context, _ interface{}) (interface{}, error) {
				err := fn(ctx, name, se, st, data)
				return nil, err
			})
			if err != nil {
				logger.Error(ctx, "TriggerSideEffect call fail.", zap.Int("SideNameType", int(data.Type)), zap.String("SideEffectName", name), zap.Any("data", data), zap.Error(err))
				return err
			}

			// 标记订单副作用状态已完成
			err = dao.MarkOrderSideEffectStatusOk(ctx, data.OrderID, data.Uid, name, int(data.Type), int64(st.OrderStatusExpireDay)*86400)
			if err != nil {
				logger.Error(ctx, "TriggerSideEffect dao.MarkOrderFlowStatusOk fail.", zap.Int("SideNameType", int(data.Type)), zap.String("SideEffectName", name), zap.Any("data", data), zap.Error(err))
				// 这里不影响主进程
			}

			return nil
		})
	}

	err = gpool.GetDefGPool().GoAndWait(fns...)
	if err != nil {
		logger.Error(ctx, "TriggerSideEffect fail.", zap.Int("SideNameType", int(data.Type)), zap.Any("data", data), zap.Error(err))
		return err
	}
	return nil
}
