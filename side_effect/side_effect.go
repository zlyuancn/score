package side_effect

import (
	"context"
	"fmt"

	"github.com/zly-app/zapp/logger"
	"go.uber.org/zap"

	"github.com/zlyuancn/score/dao"
	"github.com/zlyuancn/score/model"
)

type SideEffect interface {
	// 积分变更回调, 可能会调用多次, 业务需要自行处理幂等性(可重入)
	ScoreChange(ctx context.Context, st *model.ScoreType, flow *dao.ScoreFlowModel) error
}

var seList = make(map[string]SideEffect, 0)

// 注册副作用
func RegistrySideEffect(name string, se SideEffect) {
	l := len(seList)
	seList[name] = se
	if l == len(seList) {
		panic(fmt.Errorf("RegistrySideEffect repetition name=%s", name))
	}
}

// 取消注册副作用
func UnRegistrySideEffect(name string) {
	delete(seList, name)
}

// 触发积分变更
func TriggerScoreChange(ctx context.Context, st *model.ScoreType, flow *dao.ScoreFlowModel) error {
	for name, se := range seList {
		err := se.ScoreChange(ctx, st, flow)
		if err != nil {
			logger.Error(ctx, "TriggerScoreChange fail.", zap.String("SideEffectName", name), zap.Error(err))
			return err
		}
	}
	return nil
}
