package dao

import (
	"context"
	"database/sql"

	"github.com/zlyuancn/score/client"
)

// 积分类型表
const ScoreTypeTableName = "score_type"

type ScoreTypeModel struct {
	ID                   uint32       `json:"id"`                      // 积分类型id, 用于区分业务
	ScoreName            string       `json:"score_name"`              // 积分名, 与代码无关, 用于告诉配置人员这个积分类型是什么
	StartTime            sql.NullTime `json:"start_time"`              // 生效时间
	EndTime              sql.NullTime `json:"end_time"`                // 失效时间
	OrderStatusExpireDay uint8        `json:"order_status_expire_day"` // 订单状态保留多少天
}

// 获取所有积分类型
func GetAllScoreType(ctx context.Context) ([]*ScoreTypeModel, error) {
	const cond = `select id,score_name,start_time,end_time,order_status_expire_day from ` + ScoreTypeTableName

	var ret []*ScoreTypeModel
	err := client.ScoreTypeSqlxClient.Find(ctx, &ret, cond)
	return ret, err
}
