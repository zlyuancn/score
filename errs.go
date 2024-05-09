package score

import (
	"errors"

	"github.com/zlyuancn/score/score_type"
)

var (
	// 积分类型不存在
	ErrScoreTypeNotFound = score_type.ErrScoreTypeNotFound
	// 积分类型未生效
	ErrScoreTypeInvalid = score_type.ErrScoreTypeInvalid
	// 增加/扣除积分为0
	ErrChangeScoreValueIsZero = errors.New("change score value is zero")
	// 设置积分值小于0
	ErrSetScoreValueIsLessThanZero = errors.New("set score value is less than zero")
)

var (
	// 余额不足
	ErrInsufficientBalance = errors.New("Insufficient Balance")
	// 重入参数发生变化
	ErrReentryParamsIsChanged = errors.New("reentry params is changed")
)
