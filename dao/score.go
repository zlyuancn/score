package dao

import (
	"context"
	"errors"
	"fmt"
	"hash/crc32"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cast"

	"github.com/zly-app/component/redis"

	"github.com/zlyuancn/score/client"
	"github.com/zlyuancn/score/conf"
	"github.com/zlyuancn/score/model"
)

// 模板字符串
const (
	templateString_Uid              = "<uid>"
	templateString_Domain           = "<domain>"
	templateString_ScoreTypeID      = "<score_type_id>"
	templateString_OrderID          = "<order_id>"
	templateString_ScoreTypeIDShard = "<score_type_id_shard>"
)

// status 在redis写入的数据为  操作类型_操作状态_旧值_变更值_新的值

const (
	// 增加/扣除积分 KEYS=[积分数据key, 订单状态key]  ARGV=[增加/扣除积分值, 订单状态key有效期]
	addScoreLua = `
-- 获取订单状态
local status = redis.call('GET', KEYS[2])
-- 如果状态已写入则表示订单已完成, 直接返回状态
if status ~= false then
    return status .. '_1'
end

local changeScore = tonumber(ARGV[1])
local ex = tonumber(ARGV[2])

-- 增减积分
local nowScore = redis.call('INCRBY', KEYS[1], changeScore)
-- 检查余额不足
if nowScore < 0 then
    -- 回退
    redis.call('INCRBY', KEYS[1], -changeScore)
    -- 余额不足状态
    if changeScore > 0 then
        status = '1_2_' .. tostring(nowScore-changeScore) .. '_' .. tostring(changeScore) .. '_' .. tostring(nowScore-changeScore)
    else
        status = '2_2_' .. tostring(nowScore-changeScore) .. '_' .. tostring(-changeScore) .. '_' .. tostring(nowScore-changeScore)
    end
else
    -- 成功状态
    if changeScore > 0 then
        status = '1_1_' .. tostring(nowScore-changeScore) .. '_' .. tostring(changeScore) .. '_' .. tostring(nowScore)
    else
        status = '2_1_' .. tostring(nowScore-changeScore) .. '_' .. tostring(-changeScore) .. '_' .. tostring(nowScore)
    end
end

-- 写入状态
if ex < 1 then
    redis.call('SET', KEYS[2], status)
else
    redis.call('SET', KEYS[2], status, 'ex', ex)
end
return status .. '_0'
`

	// 重设积分 KEYS=[积分数据key, 订单状态key]  ARGV=[重设结果, 订单状态key有效期]
	resetScoreLua = `
-- 获取订单状态
local status = redis.call('GET', KEYS[2])
-- 如果状态已写入则表示订单已完成, 直接返回状态
if status ~= false then
    return status .. '_1'
end

local changeScore = tonumber(ARGV[1])
local ex = tonumber(ARGV[2])

-- 获取之前的积分
local oldScore = redis.call('GET', KEYS[1])
if oldScore == false then
    oldScore = '0'
end

-- 重设积分
redis.call('SET', KEYS[1], changeScore)
status = '3_1_' .. tostring(oldScore) .. '_' .. tostring(changeScore) .. '_' .. tostring(changeScore)

-- 写入状态
if ex < 1 then
    redis.call('SET', KEYS[2], status)
else
    redis.call('SET', KEYS[2], status, 'ex', ex)
end
return status .. '_0'
`
)

// 订单不存在
var ErrOrderNotFound = errors.New("order not found")

// 生成积分数据key
func genScoreDataKey(scoreTypeID uint32, domain string, uid string) string {
	text := conf.Conf.ScoreDataKeyFormat
	text = strings.ReplaceAll(text, templateString_ScoreTypeID, cast.ToString(scoreTypeID))
	text = strings.ReplaceAll(text, templateString_Domain, domain)
	text = strings.ReplaceAll(text, templateString_Uid, uid)
	return text
}

// 生成订单状态key
func genOrderStatusKey(uid string, orderID string) string {
	text := conf.Conf.OrderStatusKeyFormat
	text = strings.ReplaceAll(text, templateString_Uid, cast.ToString(uid))
	text = strings.ReplaceAll(text, templateString_OrderID, orderID)
	return text
}

// 生成订单序列号生成器key
func genGenOrderSeqNoKey(scoreTypeID uint32, scoreTypeIdShard int32) string {
	text := conf.Conf.GenOrderSeqNoKeyFormat
	text = strings.ReplaceAll(text, templateString_ScoreTypeID, cast.ToString(scoreTypeID))
	text = strings.ReplaceAll(text, templateString_ScoreTypeIDShard, cast.ToString(scoreTypeIdShard))
	return text
}

// 获取积分
func GetScore(ctx context.Context, scoreTypeID uint32, domain string, uid string) (int64, error) {
	key := genScoreDataKey(scoreTypeID, domain, uid)
	v, err := client.GetScoreRedisClient().Get(ctx, key).Result()
	if err == redis.Nil {
		return 0, nil
	}
	return cast.ToInt64(v), err
}

// 生成订单序列号
func GenOrderSeqNo(ctx context.Context, scoreTypeID uint32, domain string, uid string) (string, error) {
	shard := rand.Int31n(conf.Conf.GenOrderSeqNoKeyShardNum)
	key := genGenOrderSeqNoKey(scoreTypeID, shard)
	no, err := client.GetScoreRedisClient().IncrBy(ctx, key, 1).Result()
	if err != nil {
		return "", err
	}
	const orderSeqNoFormat = "%d_%d_%d_%s_%d_%s"
	uidHash := crc32.ChecksumIEEE([]byte(uid))
	uidHashDoubleHex := strconv.FormatInt(int64(uidHash), 32)

	t := time.Now().Unix()
	return fmt.Sprintf(orderSeqNoFormat, t, shard, no, uidHashDoubleHex, scoreTypeID, domain), nil
}

// 增加/扣除积分
func AddScore(ctx context.Context, orderID string, scoreTypeID uint32, domain string, uid string, score int64, statusExpireSec int64) (*model.OrderData, model.OrderStatus, error) {
	scoreDataKey := genScoreDataKey(scoreTypeID, domain, uid)
	orderStatusKey := genOrderStatusKey(uid, orderID)

	statusResult, err := client.GetScoreRedisClient().Eval(ctx, addScoreLua, []string{scoreDataKey, orderStatusKey}, score, statusExpireSec).Result()
	if err != nil {
		return nil, 0, err
	}

	return parseStatus(cast.ToString(statusResult))
}

// 重设积分
func ResetScore(ctx context.Context, orderID string, scoreTypeID uint32, domain string, uid string, resetScore int64, statusExpireSec int64) (*model.OrderData, model.OrderStatus, error) {
	scoreDataKey := genScoreDataKey(scoreTypeID, domain, uid)
	orderStatusKey := genOrderStatusKey(uid, orderID)

	statusResult, err := client.GetScoreRedisClient().Eval(ctx, resetScoreLua, []string{scoreDataKey, orderStatusKey}, resetScore, statusExpireSec).Result()
	if err != nil {
		return nil, 0, err
	}

	return parseStatus(cast.ToString(statusResult))
}

// 获取订单状态
func GetOrderStatus(ctx context.Context, orderID string, uid string) (*model.OrderData, model.OrderStatus, error) {
	orderStatusKey := genOrderStatusKey(uid, orderID)
	statusResult, err := client.GetScoreRedisClient().Get(ctx, orderStatusKey).Result()
	if err == redis.Nil {
		return nil, 0, ErrOrderNotFound
	}
	if err != nil {
		return nil, 0, err
	}
	return parseStatus(statusResult + "_0")
}

func parseStatus(statusValue string) (*model.OrderData, model.OrderStatus, error) {
	ss := strings.Split(statusValue, "_")
	if len(ss) != 6 {
		return nil, 0, fmt.Errorf("parse statusValue err. statusValue=%s", statusValue)
	}

	ret := &model.OrderData{
		OpType:      model.OpType(cast.ToInt8(ss[0])),
		OldScore:    cast.ToInt64(ss[2]),
		ChangeScore: cast.ToInt64(ss[3]),
		ResultScore: cast.ToInt64(ss[4]),
		IsReentry:   ss[5] == "1",
	}
	status := model.OrderStatus(cast.ToInt8(ss[1]))
	return ret, status, nil
}
