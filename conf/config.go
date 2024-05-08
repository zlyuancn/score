package conf

const ScoreConfigKey = "score"

const (
	defScoreRedisName           = "score"
	defScoreDataKeyFormat       = "{<uid>}:<domain>:<score_type_id>:score"
	defOrderStatusKeyFormat     = "{<uid>}:<order_id>:score_os"
	defGenOrderSeqNoKeyFormat   = "<score_type_id>:<score_type_id_shard>:score_sn"
	defGenOrderSeqNoKeyShardNum = 1000

	defScoreTypeSqlxName          = "score"
	defReloadScoreTypeIntervalSec = 60

	defScoreFlowSqlxName       = "score"
	defWriteScoreFlow          = false
	defScoreFlowTableShardNums = 2
)

var Conf = Config{
	ScoreRedisName:          defScoreRedisName,
	ScoreDataKeyFormat:      defScoreDataKeyFormat,
	OrderStatusKeyFormat:    defOrderStatusKeyFormat,
	GenOrderSeqNoKeyFormat:  defGenOrderSeqNoKeyFormat,
	ScoreTypeSqlxName:       defScoreTypeSqlxName,
	ScoreFlowSqlxName:       defScoreFlowSqlxName,
	WriteScoreFlow:          defWriteScoreFlow,
	ScoreFlowTableShardNums: defScoreFlowTableShardNums,
}

type Config struct {
	ScoreRedisName           string // 积分数据redis组件名
	ScoreDataKeyFormat       string // 积分数据key格式化字符串
	OrderStatusKeyFormat     string // 订单状态key格式化字符串
	GenOrderSeqNoKeyFormat   string // 订单号生成器key格式化字符串
	GenOrderSeqNoKeyShardNum int32  // 生成订单序列号key的分片数

	ScoreTypeSqlxName          string // 积分类型sqlx组件名
	ReloadScoreTypeIntervalSec int    // 重新加载积分类型间隔秒数

	ScoreFlowSqlxName       string // 积分流水记录sqlx组件名
	WriteScoreFlow          bool   // 是否写入积分流水
	ScoreFlowTableShardNums uint32 // 积分流水记录表分片数量
}

func (conf *Config) Check() {
	if conf.ScoreRedisName == "" {
		conf.ScoreRedisName = defScoreRedisName
	}
	if conf.ScoreDataKeyFormat == "" {
		conf.ScoreDataKeyFormat = defScoreDataKeyFormat
	}
	if conf.OrderStatusKeyFormat == "" {
		conf.OrderStatusKeyFormat = defOrderStatusKeyFormat
	}
	if conf.GenOrderSeqNoKeyFormat == "" {
		conf.GenOrderSeqNoKeyFormat = defGenOrderSeqNoKeyFormat
	}
	if conf.GenOrderSeqNoKeyShardNum < 1 {
		conf.GenOrderSeqNoKeyShardNum = defGenOrderSeqNoKeyShardNum
	}

	if conf.ScoreTypeSqlxName == "" {
		conf.ScoreTypeSqlxName = defScoreTypeSqlxName
	}
	if conf.ReloadScoreTypeIntervalSec < 1 {
		conf.ReloadScoreTypeIntervalSec = defReloadScoreTypeIntervalSec
	}

	if conf.ScoreFlowSqlxName == "" {
		conf.ScoreFlowSqlxName = defScoreFlowSqlxName
	}
	if conf.ScoreFlowTableShardNums < 1 {
		conf.ScoreFlowTableShardNums = defScoreFlowTableShardNums
	}
}
