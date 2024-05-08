package score

import (
	"context"

	"github.com/zly-app/zapp/logger"
	"go.uber.org/zap"

	"github.com/zlyuancn/score/dao"
	"github.com/zlyuancn/score/model"
	"github.com/zlyuancn/score/score_type"
)

var Score = scoreCli{}

type scoreCli struct{}

// 获取积分
func (scoreCli) GetScore(ctx context.Context, scoreTypeID uint32, domain string, uid string) (uint64, error) {
	err := score_type.CheckScoreTypeValid(ctx, scoreTypeID)
	if err != nil {
		return 0, err
	}

	score, err := dao.GetScore(ctx, scoreTypeID, domain, uid)
	if err != nil {
		logger.Log.Error(ctx, "GetScore dao.GetScore err",
			zap.Uint32("scoreTypeID", scoreTypeID),
			zap.String("domain", domain),
			zap.String("uid", uid),
			zap.Error(err),
		)
		return 0, err
	}
	return score, nil
}

// 生成订单号
func (scoreCli) GenOrderSeqNo(ctx context.Context, scoreTypeID uint32) (string, error) {
	err := score_type.CheckScoreTypeValid(ctx, scoreTypeID)
	if err != nil {
		return "", err
	}

	seqNo, err := dao.GenOrderSeqNo(ctx, scoreTypeID)
	if err != nil {
		logger.Log.Error(ctx, "GenOrderSeqNo dao.GenOrderSeqNo err",
			zap.Uint32("scoreTypeID", scoreTypeID),
			zap.Error(err),
		)
		return "", err
	}
	return seqNo, nil
}

// 增加/扣除积分
func AddScore(ctx context.Context, orderID string, scoreTypeID uint32, domain string, uid string, changeScore uint64, remark string) (*model.OrderStatusData, error) {
	if changeScore == 0 {
		logger.Log.Error(ctx, "AddScore err",
			zap.Uint32("scoreTypeID", scoreTypeID),
			zap.String("domain", domain),
			zap.String("uid", uid),
			zap.Uint64("changeScore", changeScore),
			zap.Error(ErrChangeScoreValueIsZero),
		)
		return nil, ErrChangeScoreValueIsZero
	}

	err := score_type.CheckScoreTypeValid(ctx, scoreTypeID)
	if err != nil {
		return nil, err
	}

	ret, err := dao.AddScore(ctx, orderID, scoreTypeID, domain, uid, changeScore)
	if err != nil {
		logger.Log.Error(ctx, "AddScore dao.AddScore err",
			zap.Uint32("scoreTypeID", scoreTypeID),
			zap.String("domain", domain),
			zap.String("uid", uid),
			zap.Uint64("changeScore", changeScore),
			zap.Error(ErrChangeScoreValueIsZero),
		)
		return nil, err
	}

	err = dao.WriteScoreFlow(ctx, uid, &dao.ScoreFlowModel{
		OrderID:     orderID,
		ScoreTypeID: scoreTypeID,
		Domain:      domain,
		OpType:      uint8(ret.OpType),
		OpStatus:    uint8(ret.Status),
		OldScore:    ret.OldScore,
		ChangeScore: ret.ChangeScore,
		ResultScore: ret.ResultScore,
		Uid:         uid,
		Remark:      remark,
	})
	if err != nil {
		logger.Log.Error(ctx, "AddScore dao.WriteScoreFlow err",
			zap.Uint32("scoreTypeID", scoreTypeID),
			zap.String("domain", domain),
			zap.String("uid", uid),
			zap.Uint64("changeScore", changeScore),
			zap.Error(ErrChangeScoreValueIsZero),
		)
		return nil, err
	}
	return ret, nil
}
