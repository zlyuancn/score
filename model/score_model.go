package model

import (
	"time"
)

// 积分类型
type ScoreType struct {
	ID                   uint32     // 积分类型id, 用于区分业务
	ScoreName            string     // 积分名, 与代码无关, 用于告诉配置人员这个积分类型是什么
	StartTime            *time.Time // 生效时间
	EndTime              *time.Time // 失效时间
	OrderStatusExpireDay uint8      // 订单状态保留多少天
}

// 订单数据
type OrderData struct {
	OpType      OpType // 操作类型
	OldScore    uint64 // 旧值
	ChangeScore uint64 // 变更值
	ResultScore uint64 // 新值
	IsReentry   bool   // 是否重入
}
