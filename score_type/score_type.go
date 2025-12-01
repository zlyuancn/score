package score_type

import (
	"context"
	"errors"
	"time"

	"github.com/zly-app/utils/loopload"
	"github.com/zly-app/zapp/log"
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

func StartLoopLoad() {
	loader = loopload.New("score_type", func(ctx context.Context) (map[uint32]*model.ScoreType, error) {
		if conf.Conf.ScoreTypeRedisName != "" {
			data, err := dao.GetAllScoreTypeByRedis(ctx)
			if err != nil {
				return nil, err
			}

			ret := make(map[uint32]*model.ScoreType, len(data))
			for _, d := range data {
				v := &model.ScoreType{
					ID:                        d.ID,
					ScoreName:                 d.ScoreName,
					StartTime:                 d.StartTime,
					EndTime:                   d.EndTime,
					OrderStatusExpireDay:      d.OrderStatusExpireDay,
					VerifyOrderCreateLessThan: d.VerifyOrderCreateLessThan,
				}

				if v.OrderStatusExpireDay > 0 && v.OrderStatusExpireDay < v.VerifyOrderCreateLessThan {
					v.OrderStatusExpireDay = v.VerifyOrderCreateLessThan
				}
				ret[d.ID] = v
			}
			return ret, nil
		}

		data, err := dao.GetAllScoreTypeBySqlx(ctx)
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
				v.StartTime = d.StartTime.Time.Unix()
			}
			if d.EndTime.Valid {
				v.EndTime = d.EndTime.Time.Unix()
			}

			if v.OrderStatusExpireDay > 0 && v.OrderStatusExpireDay < v.VerifyOrderCreateLessThan {
				v.OrderStatusExpireDay = v.VerifyOrderCreateLessThan
			}
			ret[d.ID] = v
		}
		return ret, nil
	}, loopload.WithReloadTime(time.Duration(conf.Conf.ReloadScoreTypeIntervalSec)*time.Second))
}

// 获取积分类型
func GetScoreType(ctx context.Context, scoreTypeID uint32) (*model.ScoreType, error) {
	st, err := getScoreType(ctx, scoreTypeID, false)
	if err != nil {
		log.Error(ctx, "GetScoreType err",
			zap.Uint32("scoreTypeID", scoreTypeID),
			zap.Error(err),
		)
	}
	return st, err
}

// 获取积分类型, 忽略积分有效期
func ForceGetScoreType(ctx context.Context, scoreTypeID uint32) (*model.ScoreType, error) {
	st, err := getScoreType(ctx, scoreTypeID, true)
	if err != nil {
		log.Error(ctx, "ForceGetScoreType err",
			zap.Uint32("scoreTypeID", scoreTypeID),
			zap.Error(err),
		)
	}
	return st, err
}

func getScoreType(ctx context.Context, scoreTypeID uint32, force bool) (*model.ScoreType, error) {
	all := loader.Get(ctx)
	st, ok := all[scoreTypeID]
	if !ok {
		return nil, ErrScoreTypeNotFound
	}

	now := int64(0)
	if !force && st.StartTime > 0 {
		now = time.Now().Unix()
		if now < st.StartTime {
			return nil, ErrScoreTypeInvalid
		}
	}

	if !force && st.EndTime > 0 {
		if now == 0 {
			now = time.Now().Unix()
		}
		if now > st.EndTime {
			return nil, ErrScoreTypeInvalid
		}
	}
	return st, nil
}
