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
				ID:                        d.ID,
				ScoreName:                 d.ScoreName,
				OrderStatusExpireDay:      d.OrderStatusExpireDay,
				VerifyOrderCreateLessThan: d.VerifyOrderCreateLessThan,
			}
			if d.StartTime.Valid {
				v.StartTime = &d.StartTime.Time
			}
			if d.EndTime.Valid {
				v.EndTime = &d.EndTime.Time
			}

			if v.OrderStatusExpireDay > 0 && v.VerifyOrderCreateLessThan > v.OrderStatusExpireDay {
				v.VerifyOrderCreateLessThan = v.OrderStatusExpireDay
			}
			ret[d.ID] = v
		}
		return ret, nil
	}, loopload.WithReloadTime(time.Duration(conf.Conf.ReloadScoreTypeIntervalSec)*time.Second))
}

// 获取积分类型
func GetScoreType(ctx context.Context, scoreTypeID uint32) (*model.ScoreType, error) {
	st, err := getScoreType(ctx, scoreTypeID)
	if err != nil {
		logger.Log.Error(ctx, "GetScoreType err",
			zap.Uint32("scoreTypeID", scoreTypeID),
			zap.Error(err),
		)
	}
	return st, err
}
func getScoreType(ctx context.Context, scoreTypeID uint32) (*model.ScoreType, error) {
	all := loader.Get(ctx)
	st, ok := all[scoreTypeID]
	if !ok {
		return nil, ErrScoreTypeNotFound
	}

	now := int64(0)
	if st.StartTime != nil {
		now = time.Now().Unix()
		if now < st.StartTime.Unix() {
			return nil, ErrScoreTypeInvalid
		}
	}

	if st.EndTime != nil {
		if now == 0 {
			now = time.Now().Unix()
		}
		if now > st.EndTime.Unix() {
			return nil, ErrScoreTypeInvalid
		}
	}
	return st, nil
}
