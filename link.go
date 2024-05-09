package score

import (
	"github.com/zlyuancn/score/model"
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
	ScoreType       = model.ScoreType
	OrderStatusData = model.OrderData
)
