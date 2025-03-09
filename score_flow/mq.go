package score_flow

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/zly-app/zapp/logger"
	"go.uber.org/zap"

	"github.com/zlyuancn/score/dao"
	"github.com/zlyuancn/score/model"
)

// 积分流水mq工具
type ScoreFlowMq interface {
	// 发送修改积分的mq信号, 延迟消费
	SendChangeScoreMqSignal(ctx context.Context, message string) error
}

var scoreFlowMq ScoreFlowMq = defScoreFlowImpl{}

type defScoreFlowImpl struct{}

func (s defScoreFlowImpl) SendChangeScoreMqSignal(ctx context.Context, message string) error {
	return nil
}

// 注入积分流水mq工具
func InjectScoreFlowMq(v ScoreFlowMq) { scoreFlowMq = v }

// 发送修改积分的mq信号
func SendChangeScoreMqSignal(ctx context.Context, cmd *model.OpCommand) error {
	message, err := sonic.MarshalString(cmd)
	if err != nil {
		logger.Error(ctx, "SendChangeScoreMqSignal call MarshalString cmd fail", zap.Any("cmd", cmd), zap.Error(err))
		return err
	}

	err = scoreFlowMq.SendChangeScoreMqSignal(ctx, message)
	if err != nil {
		logger.Error(ctx, "SendChangeScoreMqSignal call SendChangeScoreMqSignal fail", zap.String("message", message), zap.Error(err))
		return err
	}
	return nil
}

// 触发mq信号时回调
func TriggerMqSignalCallback(ctx context.Context, message string) error {
	cmd := &model.OpCommand{}
	err := sonic.UnmarshalString(message, cmd)
	if err != nil {
		logger.Error(ctx, "TriggerMqSignal call UnmarshalString cmd fail", zap.Any("message", message), zap.Error(err))
		return nil
	}

	// 获取状态
	orderData, orderStatus, err := dao.GetOrderStatus(ctx, cmd.OrderID, cmd.Uid)
	if err == dao.ErrOrderNotFound {
		logger.Error(ctx, "TriggerMqSignal call GetOrderStatus fail", zap.Any("cmd", cmd), zap.Error(err))
		// 订单不存在无需重试
		return nil
	}
	if err != nil {
		logger.Log.Error(ctx, "GetOrderStatus err",
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
	err = WriteScoreFlow(ctx, flow)
	if err != nil {
		logger.Error(ctx, "TriggerMqSignal call WriteScoreFlow fail", zap.Any("flow", flow), zap.Error(err))
		return err
	}
	return nil
}
