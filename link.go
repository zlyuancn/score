package score

import (
	"context"

	"github.com/zlyuancn/score/model"
	"github.com/zlyuancn/score/side_effect"
)

// 操作类型
type OpType = model.OpType

const (
	OpType_Add    OpType = model.OpType_Add    // 增加
	OpType_Deduct OpType = model.OpType_Deduct // 扣除
	OpType_Reset  OpType = model.OpType_Reset  // 重置
)

// 订单状态
type OrderStatus = model.OrderStatus

const (
	OrderStatus_Finish              OrderStatus = model.OrderStatus_Finish              // 完成
	OrderStatus_InsufficientBalance OrderStatus = model.OrderStatus_InsufficientBalance // 余额不足
)

type (
	ScoreType = model.ScoreType
	OrderData = model.OrderData
)

// mq工具
type MqTool = side_effect.MqTool

// 注册mq工具
func RegistryMqTool(v MqTool) {
	side_effect.RegistryMqTool(v)
}

// 触发mq回调. 触发mq信号时回调. 如果这个函数失败, 要求业务mq重试
func TriggerMqHandle(ctx context.Context, payload string) error {
	return side_effect.TriggerMqHandle(ctx, payload)
}

// 副作用
type SideEffect = side_effect.SideEffect

// 注册副作用, 重复注册同一个name会导致panic
func RegistrySideEffect(name string, se SideEffect) {
	side_effect.RegistrySideEffect(name, se)
}

// 取消注册副作用
func UnRegistrySideEffect(name string) {
	side_effect.UnRegistrySideEffect(name)
}
