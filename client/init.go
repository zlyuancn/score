package client

import (
	"github.com/zly-app/component/redis"
	"github.com/zly-app/component/sqlx"
	"github.com/zly-app/zapp/core"

	"github.com/zlyuancn/score/conf"
)

var (
	redisCreator     redis.IRedisCreator
	ScoreRedisClient redis.UniversalClient

	sqlxCreator         sqlx.ISqlx
	ScoreTypeSqlxClient sqlx.Client
	ScoreFlowSqlxClient sqlx.Client
)

func Init(app core.IApp) {
	redisCreator = redis.NewRedisCreator(app)
	ScoreRedisClient = redisCreator.GetRedis(conf.Conf.ScoreRedisName)

	sqlxCreator = sqlx.NewSqlx(app)
	ScoreTypeSqlxClient = sqlxCreator.GetSqlx(conf.Conf.ScoreTypeSqlxName)
	ScoreFlowSqlxClient = sqlxCreator.GetSqlx(conf.Conf.ScoreFlowSqlxName)
}
func Close() {
	redisCreator.Close()
	sqlxCreator.Close()
}
