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

// 副作用类型
type SideEffectType int8

const (
	SideEffectType_BeforeScoreChange SideEffectType = iota + 1 // 积分变更前
	SideEffectType_AfterScoreChange                            // 积分变更后
)

// 副作用数据
type SideEffectData struct {
	Type        SideEffectType `json:"t"`   // 副作用类型
	ScoreTypeID uint32         `json:"st"`  // 积分类型id
	Domain      string         `json:"d"`   // 积分域
	OrderID     string         `json:"oid"` // 订单id
	Uid         string         `json:"uid"` // 用户id
	Op          OpType         `json:"op"`  // 操作类型
	Score       int64          `json:"v"`   // 积分值
	Remark      string         `json:"ps"`  // 备注
}
