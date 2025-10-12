package score

import (
	"context"
	"errors"
	"fmt"
	"hash/crc32"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cast"
	"github.com/zly-app/zapp/component/gpool"
	"github.com/zly-app/zapp/pkg/utils"
	"go.uber.org/zap"

	"github.com/zly-app/zapp/logger"

	"github.com/zlyuancn/score/dao"
	"github.com/zlyuancn/score/model"
	"github.com/zlyuancn/score/score_type"
	"github.com/zlyuancn/score/side_effect"
)

var scoreApi = scoreCli{}

type scoreCli struct{}

// 获取积分
func (scoreCli) GetScore(ctx context.Context, scoreTypeID uint32, domain string, uid string) (int64, error) {
	st, err := score_type.GetScoreType(ctx, scoreTypeID)
	if err != nil {
		return 0, err
	}

	score, err := dao.GetScore(ctx, scoreTypeID, domain, uid)
	if err != nil {
		logger.Error(ctx, "GetScore dao.GetScore err",
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
func (scoreCli) GenOrderSeqNo(ctx context.Context, scoreTypeID uint32, domain string, uid string) (string, error) {
	st, err := score_type.GetScoreType(ctx, scoreTypeID)
	if err != nil {
		return "", err
	}

	seqNo, err := dao.GenOrderSeqNo(ctx, scoreTypeID, domain, uid)
	if err != nil {
		logger.Error(ctx, "GenOrderSeqNo dao.GenOrderSeqNo err",
			zap.Uint32("scoreTypeID", scoreTypeID),
			zap.String("scoreName", st.ScoreName),
			zap.String("domain", domain),
			zap.String("uid", uid),
			zap.Error(err),
		)
		return "", err
	}
	return seqNo, nil
}

// 增加积分
func (s scoreCli) AddScore(ctx context.Context, scoreTypeID uint32, domain string, uid string, orderID string, score int64, remark string) (*OrderData, error) {
	st, err := s.beforeScoreOp(ctx, model.OpType_Add, scoreTypeID, domain, uid, orderID, score, remark)
	if err != nil {
		return nil, err
	}

	// 增加积分
	data, status, err := dao.AddScore(ctx, orderID, scoreTypeID, domain, uid, score, int64(st.OrderStatusExpireDay)*86400)
	if err != nil {
		logger.Error(ctx, "AddScore dao.AddScore err",
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

	err = s.afterScoreOp(ctx, model.OpType_Add, scoreTypeID, domain, uid, orderID, score, remark, st, data, status)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// 扣除积分
func (s scoreCli) DeductScore(ctx context.Context, scoreTypeID uint32, domain string, uid string, orderID string, score int64, remark string) (*OrderData, error) {
	st, err := s.beforeScoreOp(ctx, model.OpType_Deduct, scoreTypeID, domain, uid, orderID, score, remark)
	if err != nil {
		return nil, err
	}

	// 扣除积分
	data, status, err := dao.AddScore(ctx, orderID, scoreTypeID, domain, uid, -score, int64(st.OrderStatusExpireDay)*86400)
	if err != nil {
		logger.Error(ctx, "DeductScore dao.AddScore err",
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

	err = s.afterScoreOp(ctx, model.OpType_Deduct, scoreTypeID, domain, uid, orderID, score, remark, st, data, status)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// 重设积分
func (s scoreCli) ResetScore(ctx context.Context, scoreTypeID uint32, domain string, uid string, orderID string, score int64, remark string) (*OrderData, error) {
	st, err := s.beforeScoreOp(ctx, model.OpType_Reset, scoreTypeID, domain, uid, orderID, score, remark)
	if err != nil {
		return nil, err
	}

	// 重设积分
	data, status, err := dao.ResetScore(ctx, orderID, scoreTypeID, domain, uid, score, int64(st.OrderStatusExpireDay)*86400)
	if err != nil {
		logger.Error(ctx, "ResetScore dao.ResetScore err",
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

	err = s.afterScoreOp(ctx, model.OpType_Reset, scoreTypeID, domain, uid, orderID, score, remark, st, data, status)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// 获取订单状态
func (s scoreCli) GetOrderStatus(ctx context.Context, uid string, orderID string) (*OrderData, OrderStatus, error) {
	data, status, err := dao.GetOrderStatus(ctx, orderID, uid)
	if err != nil {
		logger.Error(ctx, "GetOrderStatus err",
			zap.String("orderID", orderID),
			zap.String("uid", uid),
			zap.Error(err),
		)
		return nil, 0, err
	}
	return data, status, err
}

func (s scoreCli) beforeScoreOp(ctx context.Context, op model.OpType, scoreTypeID uint32, domain string, uid string, orderID string, score int64, remark string) (*model.ScoreType, error) {
	opName := model.GetOpName(op)
	if score < 0 {
		logger.Error(ctx, "beforeScoreOp err",
			zap.String("opName", opName),
			zap.String("orderID", orderID),
			zap.Uint32("scoreTypeID", scoreTypeID),
			zap.String("domain", domain),
			zap.String("uid", uid),
			zap.Int64("score", score),
			zap.Error(ErrChangeScoreValueIsLessThanZero),
		)
		return nil, ErrChangeScoreValueIsLessThanZero
	}

	// 检查积分类型
	st, err := score_type.GetScoreType(ctx, scoreTypeID)
	if err != nil {
		return nil, err
	}

	// 检查订单id
	err = s.verifyOrderID(orderID, scoreTypeID, domain, uid, int64(st.VerifyOrderCreateLessThan))
	if err != nil {
		logger.Error(ctx, "beforeScoreOp verifyOrderID err",
			zap.String("opName", opName),
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

	// 准备副作用
	data := &model.SideEffect{
		Type: model.SideEffectType_ScoreChange,
		ScoreChangeOpCommand: &model.OpCommand{
			Op:          op,
			ScoreTypeID: scoreTypeID,
			Domain:      domain,
			Uid:         uid,
			OrderID:     orderID,
			Score:       score,
			Remark:      remark,
		},
	}
	err = side_effect.PrepareSideEffect(ctx, data)
	if err != nil {
		logger.Error(ctx, "beforeScoreOp call mq.TriggerSendMq fail", zap.Any("cmd", data), zap.Error(err))
		return nil, err
	}

	return st, nil
}

func (s scoreCli) afterScoreOp(ctx context.Context, op model.OpType, scoreTypeID uint32, domain string, uid string, orderID string, score int64,
	remark string, st *model.ScoreType, orderData *model.OrderData, orderStatus model.OrderStatus) error {
	// 流水数据
	flow := &dao.ScoreFlowModel{
		OrderID:     orderID,
		ScoreTypeID: scoreTypeID,
		Domain:      domain,
		OpType:      uint8(orderData.OpType),
		OpStatus:    uint8(orderStatus),
		OldScore:    uint64(orderData.OldScore),
		ChangeScore: uint64(orderData.ChangeScore),
		ResultScore: uint64(orderData.ResultScore),
		Uid:         uid,
		Remark:      remark,
	}

	opName := model.GetOpName(op)

	// 检查重入时参数发生了变化
	err := s.checkReentryParamsIsChanged(orderData, op, score)
	if err != nil {
		logger.Error(ctx, "afterScoreOp checkReentryParamsIsChanged err",
			zap.String("opName", opName),
			zap.String("scoreName", st.ScoreName),
			zap.Any("flow", flow),
			zap.Error(err),
		)
		return err
	}

	// 副作用
	cloneCtx := utils.Ctx.CloneContext(ctx)
	gpool.GetDefGPool().Go(func() error {
		return side_effect.TriggerScoreChange(cloneCtx, st, flow)
	}, func(err error) {
		if err != nil {
			logger.Error(cloneCtx, "afterScoreOp call side_effect.TriggerScoreChange fail", zap.Any("flow", flow), zap.Error(err))
		}
	})

	// 检查状态
	err = s.checkStatus(orderStatus)
	if err != nil {
		logger.Error(ctx, "afterScoreOp checkStatus err",
			zap.String("opName", opName),
			zap.String("scoreName", st.ScoreName),
			zap.Int64("score", score),
			zap.Any("flow", flow),
			zap.Error(err),
		)
		return err
	}
	return nil
}

// 验证订单id
func (scoreCli) verifyOrderID(orderID string, scoreTypeID uint32, domain string, uid string, verifyOrderIDCreateLessThan int64) error {
	ss := strings.SplitN(orderID, "_", 6)
	if len(ss) != 6 {
		return errors.New("orderID invalid")
	}

	uidHash := crc32.ChecksumIEEE([]byte(uid))
	uidHashHex := strconv.FormatInt(int64(uidHash), 16)
	if ss[3] != uidHashHex {
		return errors.New("orderID not matched uid")
	}
	if ss[4] != cast.ToString(scoreTypeID) {
		return errors.New("orderID not matched scoreTypeID")
	}
	domainHash := crc32.ChecksumIEEE([]byte(domain))
	domainHashHex := strconv.FormatInt(int64(domainHash), 16)
	if ss[5] != domainHashHex {
		return errors.New("orderID not matched domain")
	}
	timestamp, err := strconv.ParseInt(ss[0], 10, 64)
	if err != nil {
		return fmt.Errorf("orderID not parsed timestamp")
	}
	if time.Now().Unix() > int64(verifyOrderIDCreateLessThan)*86400+timestamp {
		return errors.New("orderID timeout")
	}
	return nil
}

/*
	检查重入参数是否发生了变化

在订单状态key中已经包含了 uid/orderID, 而 orderID 是根据 scoreTypeID, domain 生成的, 所以无需检查这些参数
*/
func (scoreCli) checkReentryParamsIsChanged(data *model.OrderData, op model.OpType, changeScore int64) error {
	if !data.IsReentry {
		return nil
	}

	if data.OpType != op {
		return errors.New("reentry opType is changed")
	}
	if data.ChangeScore != changeScore {
		return errors.New("reentry changeScore is changed")
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
