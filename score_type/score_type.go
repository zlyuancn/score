package score_type

import (
	"context"
	"errors"
	"time"

	"github.com/zly-app/utils/loopload"
	"github.com/zly-app/zapp/core"
	"github.com/zly-app/zapp/logger"
	"go.uber.org/zap"

	"github.com/zlyuancn/score/conf"
	"github.com/zlyuancn/score/dao"

	"github.com/zlyuancn/score/model"
)

var loader *loopload.LoopLoad[map[uint32]*model.ScoreType]

var (
	// 积分类型不存在
	ErrScoreTypeNotFound = errors.New("score type not found")
	// 积分类型未生效
	ErrScoreTypeInvalid = errors.New("score type invalid")
)

func Init(app core.IApp) {
	loader = loopload.New("score_type", func(ctx context.Context) (map[uint32]*model.ScoreType, error) {
		data, err := dao.GetAllScoreType(ctx)
		if err != nil {
			return nil, err
		}

		ret := make(map[uint32]*model.ScoreType, len(data))
		for _, d := range data {
			v := &model.ScoreType{
				ID:                   d.ID,
				ScoreName:            d.ScoreName,
				OrderStatusExpireDay: d.OrderStatusExpireDay,
			}
			if d.StartTime.Valid {
				v.StartTime = &d.StartTime.Time
			}
			if d.EndTime.Valid {
				v.EndTime = &d.EndTime.Time
			}
			ret[d.ID] = v
		}
		return ret, nil
	}, loopload.WithReloadTime(time.Duration(conf.Conf.ReloadScoreTypeIntervalSec)*time.Second))
}

// 获取积分类型
func GetScoreType(ctx context.Context, scoreTypeID uint32) (*model.ScoreType, error) {
	all := loader.Get(ctx)
	ret, ok := all[scoreTypeID]
	if !ok {
		return nil, ErrScoreTypeNotFound
	}
	return ret, nil
}

// 检查积分类型有效
func CheckScoreTypeValid(ctx context.Context, scoreTypeID uint32) error {
	err := checkScoreTypeValid(ctx, scoreTypeID)
	if err != nil {
		logger.Log.Error(ctx, "CheckScoreTypeValid err",
			zap.Uint32("scoreTypeID", scoreTypeID),
			zap.Error(err),
		)
	}
	return err
}

func checkScoreTypeValid(ctx context.Context, scoreTypeID uint32) error {
	st, err := GetScoreType(ctx, scoreTypeID)
	if err != nil {
		return err
	}

	now := int64(0)
	if st.StartTime != nil {
		now = time.Now().Unix()
		if now < st.StartTime.Unix() {
			return ErrScoreTypeInvalid
		}
	}

	if st.EndTime != nil {
		if now == 0 {
			now = time.Now().Unix()
		}
		if now > st.EndTime.Unix() {
			return ErrScoreTypeInvalid
		}
	}
	return nil
}
