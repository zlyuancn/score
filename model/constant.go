package model

// 操作类型
type OpType int8

const (
	OpType_Add    OpType = 1 // 增加
	OpType_Deduct OpType = 2 // 扣除
	OpType_Reset  OpType = 3 // 重置
)

// 订单状态
type OrderStatus int8

const (
	OrderStatus_Finish              OrderStatus = 1 // 完成
	OrderStatus_InsufficientBalance OrderStatus = 2 // 余额不足
)
