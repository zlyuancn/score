package score_flow

import (
	"context"

	"github.com/zly-app/zapp/logger"
	"go.uber.org/zap"

	"github.com/zlyuancn/score/conf"
	"github.com/zlyuancn/score/dao"
	"github.com/zlyuancn/score/model"
	"github.com/zlyuancn/score/side_effect"
)

// 写入积分流水
func writeScoreFlow(ctx context.Context, st *model.ScoreType, flow *dao.ScoreFlowModel) error {
	opName := model.GetOpName(model.OpType(flow.OpType))
	err := dao.WriteScoreFlow(ctx, flow.Uid, flow)
	if err != nil {
		logger.Error(ctx, "afterScoreOp dao.writeScoreFlow fail.",
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

func (ScoreChangeSideEffect) ScoreChange(ctx context.Context, st *model.ScoreType, flow *dao.ScoreFlowModel) error {
	if !conf.Conf.WriteScoreFlow {
		return nil
	}

	// 写入流水
	err := writeScoreFlow(ctx, st, flow)
	if err != nil {
		logger.Error(ctx, "SideEffect.ScoreChange call writeScoreFlow fail.", zap.Any("flow", flow), zap.Error(err))
	}
	return nil
}
