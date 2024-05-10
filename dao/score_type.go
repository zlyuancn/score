package dao

import (
	"context"
	"database/sql"

	"github.com/zlyuancn/score/client"
)

// 积分类型表
const ScoreTypeTableName = "score_type"

type ScoreTypeModel struct {
	ID                        uint32       `db:"id"`                            // 积分类型id, 用于区分业务
	ScoreName                 string       `db:"score_name"`                    // 积分名, 与代码无关, 用于告诉配置人员这个积分类型是什么
	StartTime                 sql.NullTime `db:"start_time"`                    // 生效时间
	EndTime                   sql.NullTime `db:"end_time"`                      // 失效时间
	OrderStatusExpireDay      uint16       `db:"order_status_expire_day"`       // 订单状态保留多少天
	VerifyOrderCreateLessThan uint16       `db:"verify_order_create_less_than"` // 操作时验证订单id创建时间小于多少天
}

// 获取所有积分类型
func GetAllScoreType(ctx context.Context) ([]*ScoreTypeModel, error) {
	const cond = `select id,score_name,start_time,end_time,order_status_expire_day,verify_order_create_less_than from ` + ScoreTypeTableName

	var ret []*ScoreTypeModel
	err := client.ScoreTypeSqlxClient.Find(ctx, &ret, cond)
	return ret, err
}
