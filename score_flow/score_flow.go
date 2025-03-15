package score_flow

import (
	"context"

	"github.com/zly-app/zapp/component/gpool"
	"github.com/zly-app/zapp/logger"
	"go.uber.org/zap"

	"github.com/zlyuancn/score/dao"
	"github.com/zlyuancn/score/model"
)

// 写入积分流水
func WriteScoreFlow(ctx context.Context, st *model.ScoreType, flow *dao.ScoreFlowModel) error {
	// 获取订单流水状态
	ok, err := dao.GetOrderFlowStatus(ctx, flow.OrderID, flow.Uid)
	if err != nil {
		logger.Error(ctx, "WriteScoreFlow call GetOrderFlowStatus fail", zap.Error(err))
		// 这里不要返回err, 这里相当于没有拦截住透传到db层
	}
	if ok {
		return nil
	}

	opName := model.GetOpName(model.OpType(flow.OpType))
	err = dao.WriteScoreFlow(ctx, flow.Uid, flow)
	if err != nil {
		logger.Error(ctx, "afterScoreOp dao.WriteScoreFlow fail",
			zap.String("opName", opName),
			zap.Any("flow", flow),
			zap.Error(err),
		)
		return err
	}

	// 标记订单流水状态已落库
	gpool.GetDefGPool().Go(func() error {
		err = dao.MarkOrderFlowStatusOk(ctx, flow.OrderID, flow.Uid, int64(st.OrderStatusExpireDay)*86400)
		if err != nil {
			logger.Error(ctx, "afterScoreOp dao.MarkOrderFlowStatusOk fail", zap.Error(err))
			// 这里不影响主进程
		}
		return nil
	}, nil)

	return nil
}
