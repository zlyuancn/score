package model

import (
	"fmt"
)

// 操作类型
type OpType int8

const (
	OpType_Add    OpType = 1 // 增加
	OpType_Deduct OpType = 2 // 扣除
	OpType_Reset  OpType = 3 // 重置
)

var OpTypeName = map[OpType]string{
	OpType_Add:    "Add",
	OpType_Deduct: "Deduct",
	OpType_Reset:  "Reset",
}

func GetOpName(op OpType) string {
	n, ok := OpTypeName[op]
	if ok {
		return n
	}
	return fmt.Sprintf("Undefined op=%d", op)
}

// 订单状态
type OrderStatus int8

const (
	OrderStatus_Finish              OrderStatus = 1 // 完成
	OrderStatus_InsufficientBalance OrderStatus = 2 // 余额不足
)

// 操作指令
type OpCommand struct {
	// 操作类型
	Op OpType `json:"a"`
	// 积分类型id
	ScoreTypeID uint32 `json:"b"`
	// 积分域
	Domain string `json:"c"`
	// 用户id
	Uid string `json:"d"`
	// 订单id
	OrderID string `json:"e"`
	// 积分值
	Score int64 `json:"f"`
	// 备注
	Remark string `json:"g"`
}
