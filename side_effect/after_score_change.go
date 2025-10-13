package side_effect

import (
	"context"

	"github.com/zly-app/zapp/logger"
	"go.uber.org/zap"

	"github.com/zlyuancn/score/dao"
	"github.com/zlyuancn/score/model"
)

func afterScoreChangeHandle(ctx context.Context, st *model.ScoreType, data *model.SideEffectData) error {
	// 获取状态
	orderData, orderStatus, err := dao.GetOrderStatus(ctx, data.OrderID, data.Uid)
	if err == dao.ErrOrderNotFound {
		logger.Warn(ctx, "afterScoreChangeHandle call GetOrderStatus fail.", zap.Any("data", data), zap.Error(err))
		// 订单不存在无需重试
		return nil
	}
	if err != nil {
		logger.Error(ctx, "afterScoreChangeHandle call GetOrderStatus err", zap.Any("data", data), zap.Error(err))
		return err
	}

	// 流水数据
	flow := &dao.ScoreFlowModel{
		OrderID:     data.OrderID,
		ScoreTypeID: data.ScoreTypeID,
		Domain:      data.Domain,
		OpType:      uint8(orderData.OpType),
		OpStatus:    uint8(orderStatus),
		OldScore:    uint64(orderData.OldScore),
		ChangeScore: uint64(orderData.ChangeScore),
		ResultScore: uint64(orderData.ResultScore),
		Uid:         data.Uid,
		Remark:      data.Remark,
	}

	// 触发积分变更副作用
	err = TriggerSideEffect(ctx, data,
		func(ctx context.Context, seName string, se SideEffect, st *model.ScoreType, data *model.SideEffectData) error {
			return se.AfterScoreChange(ctx, st, data, flow)
		})
	if err != nil {
		logger.Error(ctx, "afterScoreChangeHandle call side_effect.TriggerScoreChange fail.", zap.Any("flow", flow), zap.Error(err))
		return err
	}
	return nil
}
