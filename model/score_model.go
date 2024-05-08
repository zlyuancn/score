package model

import (
	"time"
)

// 积分类型
type ScoreType struct {
	ID                   uint32
	ScoreName            string
	StartTime            *time.Time
	EndTime              *time.Time
	OrderStatusExpireDay uint8
}

// 订单状态数据
type OrderStatusData struct {
	OpType      OpType
	Status      OrderStatus
	OldScore    uint64
	ChangeScore uint64
	ResultScore uint64
}
