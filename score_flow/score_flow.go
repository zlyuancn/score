package score_flow

import (
	"context"

	"github.com/zly-app/zapp/log"
	"go.uber.org/zap"

	"github.com/zlyuancn/score/conf"
	"github.com/zlyuancn/score/dao"
	"github.com/zlyuancn/score/model"
	"github.com/zlyuancn/score/side_effect"
)

// 写入积分流水
func writeScoreFlow(ctx context.Context, flow *dao.ScoreFlowModel) error {
	opName := model.GetOpName(model.OpType(flow.OpType))
	err := dao.WriteScoreFlow(ctx, flow.Uid, flow)
	if err != nil {
		log.Error(ctx, "afterScoreOp dao.writeScoreFlow fail.",
			zap.String("opName", opName),
			zap.Any("flow", flow),
			zap.Error(err),
		)
		return err
	}
	return nil
}

type ScoreChangeSideEffect struct {
	side_effect.BaseSideEffect
}

func (ScoreChangeSideEffect) AfterScoreChange(ctx context.Context, st *model.ScoreType, data *model.SideEffectData, flow *dao.ScoreFlowModel) error {
	if !conf.Conf.WriteScoreFlow {
		return nil
	}

	// 写入流水
	err := writeScoreFlow(ctx, flow)
	if err != nil {
		log.Error(ctx, "SideEffect.AfterScoreChange call writeScoreFlow fail.", zap.Any("data", data), zap.Any("flow", flow), zap.Error(err))
		return err
	}
	return nil
}
