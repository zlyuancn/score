package client

import (
	"github.com/zly-app/component/redis"
	"github.com/zly-app/component/sqlx"

	"github.com/zlyuancn/score/conf"
)

// 获取积分 redis 客户端
func GetScoreRedisClient() redis.UniversalClient {
	return redis.GetClient(conf.Conf.ScoreRedisName)
}

// 获取积分类型 sqlx 客户端
func GetScoreTypeSqlxClient() sqlx.Client {
	return sqlx.GetClient(conf.Conf.ScoreTypeSqlxName)
}

// 获取积分类型 redis 客户端
func GetScoreTypeRedisClient() redis.UniversalClient {
	return redis.GetClient(conf.Conf.ScoreTypeRedisName)
}

// 获取积分流水 sqx 客户端
func GetScoreFlowSqlxClient() sqlx.Client {
	return sqlx.GetClient(conf.Conf.ScoreFlowSqlxName)
}
