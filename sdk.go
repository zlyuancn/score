package score

import (
	"context"
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

func (s *sdkCli) GetScore(ctx context.Context) (int64, error) {
	return scoreApi.GetScore(ctx, s.scoreTypeID, s.domain, s.uid)
}

func (s *sdkCli) GenOrderSeqNo(ctx context.Context) (string, error) {
	return scoreApi.GenOrderSeqNo(ctx, s.scoreTypeID, s.domain, s.uid)
}

func (s *sdkCli) AddScore(ctx context.Context, orderID string, score int64, remark string) (*OrderData, error) {
	return scoreApi.AddScore(ctx, s.scoreTypeID, s.domain, s.uid, orderID, score, remark)
}

func (s *sdkCli) DeductScore(ctx context.Context, orderID string, score int64, remark string) (*OrderData, error) {
	return scoreApi.DeductScore(ctx, s.scoreTypeID, s.domain, s.uid, orderID, score, remark)
}

func (s *sdkCli) ResetScore(ctx context.Context, orderID string, score int64, remark string) (*OrderData, error) {
	return scoreApi.ResetScore(ctx, s.scoreTypeID, s.domain, s.uid, orderID, score, remark)
}

func (s *sdkCli) GetOrderStatus(ctx context.Context, orderID string) (*OrderData, OrderStatus, error) {
	return scoreApi.GetOrderStatus(ctx, s.uid, orderID)
}

func NewSdk(scoreTypeID uint32, domain string, uid string) SDK {
	return &sdkCli{
		scoreTypeID: scoreTypeID,
		domain:      domain,
		uid:         uid,
	}
}
