package score

import (
	"context"

	"github.com/zly-app/zapp/filter"
)

const (
	clientType = "order"
	clientName = "sdk"
)

type SDK interface {
	// 获取积分
	GetScore(ctx context.Context) (int64, error)
	// 生成订单号
	GenOrderSeqNo(ctx context.Context) (string, error)
	// 增加积分
	AddScore(ctx context.Context, orderID string, score int64, remark string) (*OrderData, error)
	// 扣除积分
	DeductScore(ctx context.Context, orderID string, score int64, remark string) (*OrderData, error)
	// 重设积分
	ResetScore(ctx context.Context, orderID string, score int64, remark string) (*OrderData, error)
	// 获取订单状态
	GetOrderStatus(ctx context.Context, orderID string) (*OrderData, OrderStatus, error)
}

type sdkCli struct {
	scoreTypeID uint32
	domain      string
	uid         string
}

type reqBase struct {
	ScoreTypeID uint32
	Domain      string
	Uid         string
}
type rspGetScore struct {
	Score int64
}

func (s *sdkCli) GetScore(ctx context.Context) (int64, error) {
	ctx, chain := filter.GetClientFilter(ctx, clientType, clientName, "GetScore")
	r := &reqBase{
		ScoreTypeID: s.scoreTypeID,
		Domain:      s.domain,
		Uid:         s.uid,
	}
	sp := &rspGetScore{}
	err := chain.HandleInject(ctx, r, sp, func(ctx context.Context, req, rsp interface{}) error {
		sp := rsp.(*rspGetScore)
		var err error
		sp.Score, err = scoreApi.GetScore(ctx, s.scoreTypeID, s.domain, s.uid)
		return err
	})
	return sp.Score, err
}

type rspGenOrderSeqNo struct {
	SeqNo string
}

func (s *sdkCli) GenOrderSeqNo(ctx context.Context) (string, error) {
	ctx, chain := filter.GetClientFilter(ctx, clientType, clientName, "GenOrderSeqNo")
	r := &reqBase{
		ScoreTypeID: s.scoreTypeID,
		Domain:      s.domain,
		Uid:         s.uid,
	}
	sp := &rspGenOrderSeqNo{}
	err := chain.HandleInject(ctx, r, sp, func(ctx context.Context, req, rsp interface{}) error {
		sp := rsp.(*rspGenOrderSeqNo)
		var err error
		sp.SeqNo, err = scoreApi.GenOrderSeqNo(ctx, s.scoreTypeID, s.domain, s.uid)
		return err
	})
	return sp.SeqNo, err
}

type reqOSR struct {
	ScoreTypeID uint32
	Domain      string
	Uid         string
	OrderID     string
	Score       int64
	Remark      string
}
type rspD struct {
	Data *OrderData
}

func (s *sdkCli) AddScore(ctx context.Context, orderID string, score int64, remark string) (*OrderData, error) {
	ctx, chain := filter.GetClientFilter(ctx, clientType, clientName, "AddScore")
	r := &reqOSR{
		ScoreTypeID: s.scoreTypeID,
		Domain:      s.domain,
		Uid:         s.uid,
		OrderID:     orderID,
		Score:       score,
		Remark:      remark,
	}
	sp := &rspD{}
	err := chain.HandleInject(ctx, r, sp, func(ctx context.Context, req, rsp interface{}) error {
		sp := rsp.(*rspD)
		var err error
		sp.Data, err = scoreApi.AddScore(ctx, s.scoreTypeID, s.domain, s.uid, orderID, score, remark)
		return err
	})
	return sp.Data, err
}

func (s *sdkCli) DeductScore(ctx context.Context, orderID string, score int64, remark string) (*OrderData, error) {
	ctx, chain := filter.GetClientFilter(ctx, clientType, clientName, "DeductScore")
	r := &reqOSR{
		ScoreTypeID: s.scoreTypeID,
		Domain:      s.domain,
		Uid:         s.uid,
		OrderID:     orderID,
		Score:       score,
		Remark:      remark,
	}
	sp := &rspD{}
	err := chain.HandleInject(ctx, r, sp, func(ctx context.Context, req, rsp interface{}) error {
		sp := rsp.(*rspD)
		var err error
		sp.Data, err = scoreApi.DeductScore(ctx, s.scoreTypeID, s.domain, s.uid, orderID, score, remark)
		return err
	})
	return sp.Data, err
}

func (s *sdkCli) ResetScore(ctx context.Context, orderID string, score int64, remark string) (*OrderData, error) {
	ctx, chain := filter.GetClientFilter(ctx, clientType, clientName, "ResetScore")
	r := &reqOSR{
		ScoreTypeID: s.scoreTypeID,
		Domain:      s.domain,
		Uid:         s.uid,
		OrderID:     orderID,
		Score:       score,
		Remark:      remark,
	}
	sp := &rspD{}
	err := chain.HandleInject(ctx, r, sp, func(ctx context.Context, req, rsp interface{}) error {
		sp := rsp.(*rspD)
		var err error
		sp.Data, err = scoreApi.ResetScore(ctx, s.scoreTypeID, s.domain, s.uid, orderID, score, remark)
		return err
	})
	return sp.Data, err
}

type reqO struct {
	ScoreTypeID uint32
	Domain      string
	Uid         string
	OrderID     string
}
type rspDS struct {
	Data   *OrderData
	Status OrderStatus
}

func (s *sdkCli) GetOrderStatus(ctx context.Context, orderID string) (*OrderData, OrderStatus, error) {
	ctx, chain := filter.GetClientFilter(ctx, clientType, clientName, "GetOrderStatus")
	r := &reqO{
		ScoreTypeID: s.scoreTypeID,
		Domain:      s.domain,
		Uid:         s.uid,
		OrderID:     orderID,
	}
	sp := &rspDS{}
	err := chain.HandleInject(ctx, r, sp, func(ctx context.Context, req, rsp interface{}) error {
		sp := rsp.(*rspDS)
		var err error
		sp.Data, sp.Status, err = scoreApi.GetOrderStatus(ctx, s.uid, orderID)
		return err
	})
	return sp.Data, sp.Status, err
}

func NewSdk(scoreTypeID uint32, domain string, uid string) SDK {
	return &sdkCli{
		scoreTypeID: scoreTypeID,
		domain:      domain,
		uid:         uid,
	}
}
