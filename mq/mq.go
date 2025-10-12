package mq

import (
	"context"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/zly-app/zapp/logger"
	"go.uber.org/zap"

	"github.com/zlyuancn/score/dao"
	"github.com/zlyuancn/score/model"
	"github.com/zlyuancn/score/score_type"
	"github.com/zlyuancn/score/side_effect"
)

// mq工具
type MqTool interface {
	// 发送mq数据, 要求必须延迟10秒以上消费. 消息被消费时需要调用 TriggerMqHandle
	Send(ctx context.Context, payload string) error
}

var mqTool MqTool = BaseMqTool{}

type BaseMqTool struct{}

func (s BaseMqTool) Send(ctx context.Context, payload string) error {
	return nil
}

// 注册mq工具
func RegistryMqTool(v MqTool) { mqTool = v }

// 触发发送mq数据
func TriggerSendMq(ctx context.Context, data *model.MqData) error {
	payload, err := sonic.MarshalString(data)
	if err != nil {
		logger.Error(ctx, "TriggerSendMq call MarshalString data fail", zap.Any("data", data), zap.Error(err))
		return err
	}

	err = mqTool.Send(ctx, payload)
	if err != nil {
		logger.Error(ctx, "TriggerSendMq call mqTool.Send fail", zap.String("payload", payload), zap.Error(err))
		return err
	}
	return nil
}

// 触发mq回调. 触发mq信号时回调. 如果这个函数失败, 要求业务mq重试
func TriggerMqHandle(ctx context.Context, payload string) error {
	data := &model.MqData{}
	err := sonic.UnmarshalString(payload, data)
	if err != nil {
		logger.Error(ctx, "TriggerMqHandle call UnmarshalString data fail", zap.Any("payload", payload), zap.Error(err))
		return nil
	}

	switch data.Type {
	case model.MqDataType_ScoreChange:
		return ScoreChangeHandle(ctx, data)
	}
	return fmt.Errorf("TriggerMqHandle got not supported type=%d", data.Type)
}

func ScoreChangeHandle(ctx context.Context, data *model.MqData) error {
	cmd := data.ScoreChangeOpCommand

	// 检查积分类型
	st, err := score_type.ForceGetScoreType(ctx, cmd.ScoreTypeID)
	if err != nil {
		return err
	}

	// 获取状态
	orderData, orderStatus, err := dao.GetOrderStatus(ctx, cmd.OrderID, cmd.Uid)
	if err == dao.ErrOrderNotFound {
		logger.Warn(ctx, "ScoreChangeHandle call GetOrderStatus fail", zap.Any("cmd", cmd), zap.Error(err))
		// 订单不存在无需重试
		return nil
	}
	if err != nil {
		logger.Error(ctx, "ScoreChangeHandle call GetOrderStatus err",
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

	err = side_effect.TriggerScoreChange(ctx, st, flow)
	if err != nil {
		logger.Error(ctx, "ScoreChangeHandle call side_effect.TriggerScoreChange fail", zap.Any("flow", flow), zap.Error(err))
		return err
	}
	return nil
}
