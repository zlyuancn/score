package side_effect

import (
	"context"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/zly-app/zapp/logger"
	"go.uber.org/zap"

	"github.com/zlyuancn/score/dao"
	"github.com/zlyuancn/score/model"
)

// 副作用, 必须继承 BaseSideEffect
type SideEffect interface {
	abstract()
	// 积分变更回调, 可能会调用多次, 业务需要自行处理幂等性(可重入)
	ScoreChange(ctx context.Context, st *model.ScoreType, flow *dao.ScoreFlowModel) error
}

type BaseSideEffect struct{}

func (BaseSideEffect) abstract() {}

var seMap = make(map[model.SideEffectType]map[string]SideEffect, 0)

// 注册副作用, 重复注册同一个name会导致panic
func RegistrySideEffect(t model.SideEffectType, name string, se SideEffect) {
	seList, ok := seMap[t]
	if !ok {
		seMap[t] = map[string]SideEffect{
			name: se,
		}
		return
	}

	l := len(seList)
	seList[name] = se
	if l == len(seList) {
		panic(fmt.Errorf("RegistrySideEffect repetition name=%s type=%d", name, int(t)))
	}
}

// 取消注册副作用
func UnRegistrySideEffect(t model.SideEffectType, name string) {
	seList, ok := seMap[t]
	if ok {
		delete(seList, name)
	}
}

var sideEffectTypeResolver = map[model.SideEffectType]func(ctx context.Context, data *model.SideEffect) error{
	model.SideEffectType_ScoreChange: scoreChangeHandle,
}

// 准备副作用
func PrepareSideEffect(ctx context.Context, data *model.SideEffect) error {
	payload, err := sonic.MarshalString(data)
	if err != nil {
		logger.Error(ctx, "PrepareSideEffect call MarshalString data fail", zap.Any("data", data), zap.Error(err))
		return err
	}

	err = mqTool.Send(ctx, payload)
	if err != nil {
		logger.Error(ctx, "PrepareSideEffect call mqTool.Send fail", zap.String("payload", payload), zap.Error(err))
		return err
	}
	return nil
}

// 副作用补偿
func compensationSideEffect(ctx context.Context, payload string) error {
	data := &model.SideEffect{}
	err := sonic.UnmarshalString(payload, data)
	if err != nil {
		logger.Error(ctx, "compensationSideEffect call UnmarshalString data fail", zap.Any("payload", payload), zap.Error(err))
		return nil
	}

	r, ok := sideEffectTypeResolver[data.Type]
	if ok {
		return r(ctx, data)
	}
	return fmt.Errorf("compensationSideEffect got not supported type=%d", data.Type)
}
