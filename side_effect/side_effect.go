package side_effect

import (
	"context"

	"github.com/zly-app/zapp/logger"
	"go.uber.org/zap"

	"github.com/zlyuancn/score/dao"
	"github.com/zlyuancn/score/model"
)

// 副作用, 必须继承 BaseSideEffect, 副作用的每个方法都可能会调用多次, 业务需要自行处理幂等性(可重入)
type SideEffect interface {
	abstract()
	// 积分变更前, 如果返回err, 则积分变更会失败
	BeforeScoreChange(ctx context.Context, st *model.ScoreType, data *model.SideEffectData) error
	// 积分变更后
	AfterScoreChange(ctx context.Context, st *model.ScoreType, data *model.SideEffectData, flow *dao.ScoreFlowModel) error
}

var _ SideEffect = BaseSideEffect{}

type BaseSideEffect struct{}

func (BaseSideEffect) abstract() {}

func (e BaseSideEffect) BeforeScoreChange(ctx context.Context, st *model.ScoreType, data *model.SideEffectData) error {
	return nil
}

func (e BaseSideEffect) AfterScoreChange(ctx context.Context, st *model.ScoreType, data *model.SideEffectData, flow *dao.ScoreFlowModel) error {
	return nil
}

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
