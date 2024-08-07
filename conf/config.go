package conf

const ScoreConfigKey = "score"

const (
	defScoreRedisName           = "score"
	defScoreDataKeyFormat       = "{<uid>}:<domain>:<score_type_id>:score"
	defOrderStatusKeyFormat     = "{<uid>}:<order_id>:score_os"
	defGenOrderSeqNoKeyFormat   = "<score_type_id>:<score_type_id_shard>:score_sn"
	defGenOrderSeqNoKeyShardNum = 1000

	defScoreTypeRedisName         = "score"
	defScoreTypeRedisKey          = "score:score_type"
	defReloadScoreTypeIntervalSec = 60

	defScoreFlowSqlxName       = "score"
	defWriteScoreFlow          = false
	defScoreFlowTableShardNums = 2
)

var Conf = Config{
	ScoreRedisName:           defScoreRedisName,
	ScoreDataKeyFormat:       defScoreDataKeyFormat,
	OrderStatusKeyFormat:     defOrderStatusKeyFormat,
	GenOrderSeqNoKeyFormat:   defGenOrderSeqNoKeyFormat,
	GenOrderSeqNoKeyShardNum: defGenOrderSeqNoKeyShardNum,

	ScoreTypeRedisName:         defScoreTypeRedisName,
	ScoreTypeRedisKey:          defScoreTypeRedisKey,
	ScoreTypeSqlxName:          "",
	ReloadScoreTypeIntervalSec: defReloadScoreTypeIntervalSec,

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

	ScoreTypeRedisName         string // 积分类型redis组件名
	ScoreTypeRedisKey          string // 积分类型从redis加载的 hash map key名
	ScoreTypeSqlxName          string // 积分类型sqlx组件名, 如果配置了 ScoreTypeRedisName 则仅从redis加载积分类型
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

	if conf.ScoreTypeRedisName == "" && conf.ScoreTypeSqlxName == "" {
		conf.ScoreTypeRedisName = defScoreTypeRedisName
	}
	if conf.ScoreTypeRedisKey == "" {
		conf.ScoreTypeRedisKey = defScoreTypeRedisKey
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
