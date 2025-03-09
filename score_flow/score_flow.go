package score_flow

import (
	"context"

	"github.com/zly-app/zapp/logger"
	"go.uber.org/zap"

	"github.com/zlyuancn/score/dao"
	"github.com/zlyuancn/score/model"
)

// 写入积分流水
func WriteScoreFlow(ctx context.Context, flow *dao.ScoreFlowModel) error {
	opName := model.GetOpName(model.OpType(flow.OpType))
	err := dao.WriteScoreFlow(ctx, flow.Uid, flow)
	if err != nil {
		logger.Log.Error(ctx, "afterScoreOp dao.WriteScoreFlow fail",
			zap.String("opName", opName),
			zap.Any("flow", flow),
			zap.Error(err),
		)
		return err
	}
	return nil
}
