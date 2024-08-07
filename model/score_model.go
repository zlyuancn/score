package model

// 积分类型
type ScoreType struct {
	ID                        uint32 // 积分类型id, 用于区分业务
	ScoreName                 string // 积分名, 与代码无关, 用于告诉配置人员这个积分类型是什么
	StartTime                 int64  // 生效时间
	EndTime                   int64  // 失效时间
	OrderStatusExpireDay      uint16 // 订单状态保留多少天
	VerifyOrderCreateLessThan uint16 // 操作时验证订单id创建时间小于多少天
}

// 订单数据
type OrderData struct {
	OpType      OpType // 操作类型
	OldScore    int64  // 旧值
	ChangeScore int64  // 变更值
	ResultScore int64  // 新值
	IsReentry   bool   // 是否重入
}
