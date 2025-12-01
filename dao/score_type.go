package dao

import (
	"context"
	"database/sql"

	"github.com/bytedance/sonic"
	"github.com/spf13/cast"
	"github.com/zly-app/zapp/log"
	"go.uber.org/zap"

	"github.com/zlyuancn/score/client"
	"github.com/zlyuancn/score/conf"
)

type ScoreTypeRedisModel struct {
	ID                        uint32 `json:"-"`                             // 积分类型id, 用于区分业务
	ScoreName                 string `json:"score_name"`                    // 积分名, 与代码无关, 用于告诉配置人员这个积分类型是什么业务
	StartTime                 int64  `json:"start_time"`                    // 生效时间, 0 表示不限制
	EndTime                   int64  `json:"end_time"`                      // 失效时间, 0 表示不限制
	OrderStatusExpireDay      uint16 `json:"order_status_expire_day"`       // 订单状态保留多少天
	VerifyOrderCreateLessThan uint16 `json:"verify_order_create_less_than"` // 操作时验证订单id创建时间小于多少天
}

// 获取所有积分类型
func GetAllScoreTypeByRedis(ctx context.Context) ([]*ScoreTypeRedisModel, error) {
	rdb, err := client.GetScoreTypeRedisClient()
	if err != nil {
		return nil, err
	}

	v, err := rdb.HGetAll(ctx, conf.Conf.ScoreTypeRedisKey).Result()
	if err != nil {
		return nil, err
	}

	ret := make([]*ScoreTypeRedisModel, 0, len(v))
	for k, text := range v {
		id, err := cast.ToUint32E(k)
		if err != nil {
			log.Error(ctx, "can't parse score type by redis hash map field", zap.String("field", k), zap.String("value", text), zap.Error(err))
			return nil, err
		}
		s := ScoreTypeRedisModel{}
		err = sonic.UnmarshalString(text, &s)
		if err != nil {
			log.Error(ctx, "can't parse score type conf by redis hash map value", zap.String("field", k), zap.String("value", text), zap.Error(err))
			return nil, err
		}
		s.ID = id
		ret = append(ret, &s)
	}
	return ret, nil
}

type ScoreTypeSqlxModel struct {
	ID                        uint32       `db:"id"`                            // 积分类型id, 用于区分业务
	ScoreName                 string       `db:"score_name"`                    // 积分名, 与代码无关, 用于告诉配置人员这个积分类型是什么业务
	StartTime                 sql.NullTime `db:"start_time"`                    // 生效时间
	EndTime                   sql.NullTime `db:"end_time"`                      // 失效时间
	OrderStatusExpireDay      uint16       `db:"order_status_expire_day"`       // 订单状态保留多少天
	VerifyOrderCreateLessThan uint16       `db:"verify_order_create_less_than"` // 操作时验证订单id创建时间小于多少天
}

// 获取所有积分类型
func GetAllScoreTypeBySqlx(ctx context.Context) ([]*ScoreTypeSqlxModel, error) {
	const cond = `select id,score_name,start_time,end_time,order_status_expire_day,verify_order_create_less_than from score_type`

	var ret []*ScoreTypeSqlxModel
	err := client.GetScoreTypeSqlxClient().Find(ctx, &ret, cond)
	return ret, err
}
