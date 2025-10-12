package side_effect

import (
	"context"

	"github.com/zly-app/zapp/logger"
	"go.uber.org/zap"

	"github.com/zlyuancn/score/dao"
	"github.com/zlyuancn/score/model"
	"github.com/zlyuancn/score/score_type"
)

func scoreChangeHandle(ctx context.Context, data *model.SideEffect) error {
	cmd := data.ScoreChangeOpCommand

	// 检查积分类型
	st, err := score_type.ForceGetScoreType(ctx, cmd.ScoreTypeID)
	if err != nil {
		return err
	}

	// 获取状态
	orderData, orderStatus, err := dao.GetOrderStatus(ctx, cmd.OrderID, cmd.Uid)
	if err == dao.ErrOrderNotFound {
		logger.Warn(ctx, "scoreChangeHandle call GetOrderStatus fail", zap.Any("cmd", cmd), zap.Error(err))
		// 订单不存在无需重试
		return nil
	}
	if err != nil {
		logger.Error(ctx, "scoreChangeHandle call GetOrderStatus err",
			zap.String("orderID", cmd.OrderID),
			zap.String("uid", cmd.Uid),
			zap.Error(err),
		)
		return err
	}

	// 流水数据
	flow := &dao.ScoreFlowModel{
		OrderID:     cmd.OrderID,
		ScoreTypeID: cmd.ScoreTypeID,
		Domain:      cmd.Domain,
		OpType:      uint8(orderData.OpType),
		OpStatus:    uint8(orderStatus),
		OldScore:    uint64(orderData.OldScore),
		ChangeScore: uint64(orderData.ChangeScore),
		ResultScore: uint64(orderData.ResultScore),
		Uid:         cmd.Uid,
		Remark:      cmd.Remark,
	}

	err = TriggerScoreChange(ctx, st, flow)
	if err != nil {
		logger.Error(ctx, "scoreChangeHandle call side_effect.TriggerScoreChange fail", zap.Any("flow", flow), zap.Error(err))
		return err
	}
	return nil
}

// 立即触发积分变更副作用
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
