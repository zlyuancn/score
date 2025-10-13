package side_effect

import (
	"context"

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
		logger.Panic("RegistrySideEffect repetition name", zap.Int("SideTypeType", int(t)), zap.String("Name", name))
	}
}

// 取消注册副作用
func UnRegistrySideEffect(t model.SideEffectType, name string) {
	seList, ok := seMap[t]
	if ok {
		delete(seList, name)
	}
}
