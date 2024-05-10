package score

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cast"
	"go.uber.org/zap"

	"github.com/zly-app/zapp/logger"

	"github.com/zlyuancn/score/conf"
	"github.com/zlyuancn/score/dao"
	"github.com/zlyuancn/score/model"
	"github.com/zlyuancn/score/score_type"
)

var Score = scoreCli{}

type scoreCli struct{}

// 获取积分
func (scoreCli) GetScore(ctx context.Context, scoreTypeID uint32, domain string, uid string) (int64, error) {
	st, err := score_type.GetScoreType(ctx, scoreTypeID)
	if err != nil {
		return 0, err
	}

	score, err := dao.GetScore(ctx, scoreTypeID, domain, uid)
	if err != nil {
		logger.Log.Error(ctx, "GetScore dao.GetScore err",
			zap.Uint32("scoreTypeID", scoreTypeID),
			zap.String("scoreName", st.ScoreName),
			zap.String("domain", domain),
			zap.String("uid", uid),
			zap.Error(err),
		)
		return 0, err
	}
	return score, nil
}

// 生成订单号
func (scoreCli) GenOrderSeqNo(ctx context.Context, scoreTypeID uint32, domain string) (string, error) {
	st, err := score_type.GetScoreType(ctx, scoreTypeID)
	if err != nil {
		return "", err
	}

	seqNo, err := dao.GenOrderSeqNo(ctx, scoreTypeID, domain)
	if err != nil {
		logger.Log.Error(ctx, "GenOrderSeqNo dao.GenOrderSeqNo err",
			zap.Uint32("scoreTypeID", scoreTypeID),
			zap.String("scoreName", st.ScoreName),
			zap.String("domain", domain),
			zap.Error(err),
		)
		return "", err
	}
	return seqNo, nil
}

// 增加/扣除积分
func (s scoreCli) AddScore(ctx context.Context, orderID string, scoreTypeID uint32, domain string, uid string, score int64, remark string) (*model.OrderData, error) {
	if score == 0 {
		logger.Log.Error(ctx, "AddScore err",
			zap.String("orderID", orderID),
			zap.Uint32("scoreTypeID", scoreTypeID),
			zap.String("domain", domain),
			zap.String("uid", uid),
			zap.Int64("score", score),
			zap.Error(ErrChangeScoreValueIsZero),
		)
		return nil, ErrChangeScoreValueIsZero
	}

	// 检查积分类型
	st, err := score_type.GetScoreType(ctx, scoreTypeID)
	if err != nil {
		return nil, err
	}

	// 检查订单id
	err = s.verifyOrderID(orderID, scoreTypeID, domain)
	if err != nil {
		logger.Log.Error(ctx, "AddScore verifyOrderID err",
			zap.String("orderID", orderID),
			zap.Uint32("scoreTypeID", scoreTypeID),
			zap.String("scoreName", st.ScoreName),
			zap.String("domain", domain),
			zap.String("uid", uid),
			zap.Int64("score", score),
			zap.Error(err),
		)
		return nil, err
	}

	// 增减积分
	data, status, err := dao.AddScore(ctx, orderID, scoreTypeID, domain, uid, score, int64(st.OrderStatusExpireDay)*86400)
	if err != nil {
		logger.Log.Error(ctx, "AddScore dao.AddScore err",
			zap.String("orderID", orderID),
			zap.Uint32("scoreTypeID", scoreTypeID),
			zap.String("scoreName", st.ScoreName),
			zap.String("domain", domain),
			zap.String("uid", uid),
			zap.Int64("score", score),
			zap.Error(err),
		)
		return nil, err
	}

	// 写入流水
	flow := &dao.ScoreFlowModel{
		OrderID:     orderID,
		ScoreTypeID: scoreTypeID,
		Domain:      domain,
		OpType:      uint8(data.OpType),
		OpStatus:    uint8(status),
		OldScore:    uint64(data.OldScore),
		ChangeScore: uint64(data.ChangeScore),
		ResultScore: uint64(data.ResultScore),
		Uid:         uid,
		Remark:      remark,
	}
	if conf.Conf.WriteScoreFlow {
		err = dao.WriteScoreFlow(ctx, uid, flow)
		if err != nil {
			return nil, err
		}
	}

	// 检查状态
	err = s.checkStatus(status)
	if err != nil {
		logger.Log.Error(ctx, "AddScore checkStatus err",
			zap.String("scoreName", st.ScoreName),
			zap.Any("flow", flow),
			zap.Error(err),
		)
		return nil, err
	}

	// 检查重入时参数发生了变化
	op := model.OpType_Add
	changeScore := score
	if score < 0 {
		op = model.OpType_Deduct
		changeScore = -score
	}
	if s.checkReentryParamsIsChanged(data, op, changeScore) {
		logger.Log.Error(ctx, "AddScore checkReentryParamsIsChanged err",
			zap.String("scoreName", st.ScoreName),
			zap.Int64("score", score),
			zap.Any("flow", flow),
			zap.Error(ErrReentryParamsIsChanged),
		)
		return nil, ErrReentryParamsIsChanged
	}

	return data, nil
}

// 重设积分
func (s scoreCli) ResetScore(ctx context.Context, orderID string, scoreTypeID uint32, domain string, uid string, score int64, remark string) (*model.OrderData, error) {
	if score < 0 {
		logger.Log.Error(ctx, "ResetScore err",
			zap.String("orderID", orderID),
			zap.Uint32("scoreTypeID", scoreTypeID),
			zap.String("domain", domain),
			zap.String("uid", uid),
			zap.Int64("score", score),
			zap.Error(ErrSetScoreValueIsLessThanZero),
		)
		return nil, ErrSetScoreValueIsLessThanZero
	}

	// 检查积分类型
	st, err := score_type.GetScoreType(ctx, scoreTypeID)
	if err != nil {
		return nil, err
	}

	// 检查订单id
	err = s.verifyOrderID(orderID, scoreTypeID, domain)
	if err != nil {
		logger.Log.Error(ctx, "ResetScore verifyOrderID err",
			zap.String("orderID", orderID),
			zap.Uint32("scoreTypeID", scoreTypeID),
			zap.String("scoreName", st.ScoreName),
			zap.String("domain", domain),
			zap.String("uid", uid),
			zap.Int64("score", score),
			zap.Error(err),
		)
		return nil, err
	}

	// 重设积分
	data, status, err := dao.ResetScore(ctx, orderID, scoreTypeID, domain, uid, score, int64(st.OrderStatusExpireDay)*86400)
	if err != nil {
		logger.Log.Error(ctx, "ResetScore dao.ResetScore err",
			zap.String("orderID", orderID),
			zap.Uint32("scoreTypeID", scoreTypeID),
			zap.String("scoreName", st.ScoreName),
			zap.String("domain", domain),
			zap.String("uid", uid),
			zap.Int64("score", score),
			zap.Error(err),
		)
		return nil, err
	}

	// 写入流水
	flow := &dao.ScoreFlowModel{
		OrderID:     orderID,
		ScoreTypeID: scoreTypeID,
		Domain:      domain,
		OpType:      uint8(data.OpType),
		OpStatus:    uint8(status),
		OldScore:    uint64(data.OldScore),
		ChangeScore: uint64(data.ChangeScore),
		ResultScore: uint64(data.ResultScore),
		Uid:         uid,
		Remark:      remark,
	}
	if conf.Conf.WriteScoreFlow {
		err = dao.WriteScoreFlow(ctx, uid, flow)
		if err != nil {
			return nil, err
		}
	}

	// 检查状态
	err = s.checkStatus(status)
	if err != nil {
		logger.Log.Error(ctx, "ResetScore checkStatus err",
			zap.String("scoreName", st.ScoreName),
			zap.Int64("score", score),
			zap.Any("flow", flow),
			zap.Error(err),
		)
		return nil, err
	}

	// 检查重入时参数发生了变化
	if s.checkReentryParamsIsChanged(data, model.OpType_Reset, score) {
		logger.Log.Error(ctx, "ResetScore checkReentryParamsIsChanged err",
			zap.String("scoreName", st.ScoreName),
			zap.Any("flow", flow),
			zap.Error(ErrReentryParamsIsChanged),
		)
		return nil, ErrReentryParamsIsChanged
	}

	return data, nil
}

// 验证订单id
func (scoreCli) verifyOrderID(orderID string, scoreTypeID uint32, domain string) error {
	ss := strings.SplitN(orderID, "_", 5)
	if len(ss) != 5 {
		return errors.New("orderID invalid")
	}

	if ss[3] != cast.ToString(scoreTypeID) {
		return errors.New("orderID not matched scoreTypeID")
	}
	if ss[4] != cast.ToString(domain) {
		return errors.New("orderID not matched domain")
	}
	if time.Now().Unix() > int64(conf.Conf.VerifyOrderIDCreateLessThan)*86400+cast.ToInt64(ss[2])/1000 {
		return errors.New("orderID timeout")
	}
	return nil
}

// 检查状态
func (scoreCli) checkStatus(status model.OrderStatus) error {
	switch status {
	case model.OrderStatus_Finish:
		return nil
	case model.OrderStatus_InsufficientBalance:
		return ErrInsufficientBalance
	}
	return fmt.Errorf("undefined status=%d", status)
}

/*
	检查重入参数是否发生了变化

在订单状态key中已经包含了 uid/orderID, 而 orderID 是根据 scoreTypeID, domain 生成的, 所以无需检查这些参数
*/
func (scoreCli) checkReentryParamsIsChanged(data *model.OrderData, op model.OpType, changeScore int64) bool {
	if !data.IsReentry {
		return false
	}

	if data.OpType != op {
		return true
	}
	if data.ChangeScore != changeScore {
		return true
	}
	return false
}
