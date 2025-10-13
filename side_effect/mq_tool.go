package side_effect

import (
	"context"
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

// 触发mq回调. 触发mq信号时回调. 如果这个函数失败, 要求业务mq重试
func TriggerMqHandle(ctx context.Context, payload string) error {
	return compensationSideEffect(ctx, payload)
}
