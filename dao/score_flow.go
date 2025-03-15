package dao

import (
	"context"
	"errors"
	"hash/crc32"

	"github.com/didi/gendry/builder"
	"github.com/spf13/cast"
	"github.com/zly-app/zapp/logger"
	"go.uber.org/zap"

	"github.com/zlyuancn/score/client"
	"github.com/zlyuancn/score/conf"
)

// 积分流水表
const ScoreFlowTableName = "score_flow_"

type ScoreFlowModel struct {
	OrderID     string `db:"oid"`           // 订单id
	ScoreTypeID uint32 `db:"score_type_id"` // 积分类型id
	Domain      string `db:"domain"`        // 域

	OpType   uint8 `db:"o_type"`   // 操作类型. 1=增加, 2=扣除, 3=重置
	OpStatus uint8 `db:"o_status"` // 操作状态. 1=成功, 2=余额不足

	OldScore    uint64 `db:"old_score"`    // 原始积分
	ChangeScore uint64 `db:"change_score"` // 变更积分
	ResultScore uint64 `db:"result_score"` // 结果积分

	Uid    string `db:"uid"`    // 唯一标识一个用户
	Remark string `db:"remark"` // 备注
}

// 写入积分流水
func WriteScoreFlow(ctx context.Context, uid string, v *ScoreFlowModel) error {
	if v == nil {
		return errors.New("CreateOneModel v is empty")
	}

	shardID := crc32.ChecksumIEEE([]byte(uid)) % conf.Conf.ScoreFlowTableShardNums
	tabName := ScoreFlowTableName + cast.ToString(shardID)

	var data []map[string]interface{}
	data = append(data, map[string]interface{}{
		"oid":           v.OrderID,
		"score_type_id": v.ScoreTypeID,
		"domain":        v.Domain,

		"o_type":   v.OpType,
		"o_status": v.OpStatus,

		"old_score":    v.OldScore,
		"change_score": v.ChangeScore,
		"result_score": v.ResultScore,

		"uid":    v.Uid,
		"remark": v.Remark,
	})
	cond, vals, err := builder.BuildInsertIgnore(tabName, data)
	if err != nil {
		logger.Error(ctx, "score CreateOneModel BuildSelect err",
			zap.Any("data", data),
			zap.Error(err),
		)
		return err
	}

	result, err := client.GetScoreFlowSqlxClient().Exec(ctx, cond, vals...)
	if err != nil {
		logger.Error(ctx, "score CreateOneModel err",
			zap.String("cond", cond),
			zap.Any("vals", vals),
			zap.Error(err),
		)
		return err
	}
	_, err = result.LastInsertId()
	return err
}
