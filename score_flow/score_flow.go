package score_flow

import (
	"context"

	"github.com/zly-app/zapp/component/gpool"
	"github.com/zly-app/zapp/logger"
	"github.com/zly-app/zapp/pkg/utils"
	"go.uber.org/zap"

	"github.com/zlyuancn/score/conf"
	"github.com/zlyuancn/score/dao"
	"github.com/zlyuancn/score/model"
	"github.com/zlyuancn/score/side_effect"
)

// 写入积分流水
func writeScoreFlow(ctx context.Context, st *model.ScoreType, flow *dao.ScoreFlowModel) error {
	// 获取订单流水落库状态
	ok, err := dao.GetOrderFlowStatus(ctx, flow.OrderID, flow.Uid)
	if err != nil {
		logger.Error(ctx, "writeScoreFlow call GetOrderFlowStatus fail", zap.Error(err))
		// 这里不要返回err, 这里相当于没有拦截住透传到db层
	}
	if ok { // 已经落库成功则忽略
		return nil
	}

	opName := model.GetOpName(model.OpType(flow.OpType))
	err = dao.WriteScoreFlow(ctx, flow.Uid, flow)
	if err != nil {
		logger.Error(ctx, "afterScoreOp dao.writeScoreFlow fail",
			zap.String("opName", opName),
			zap.Any("flow", flow),
			zap.Error(err),
		)
		return err
	}

	// 标记订单流水状态已落库
	cloneCtx := utils.Ctx.CloneContext(ctx)
	gpool.GetDefGPool().Go(func() error {
		err = dao.MarkOrderFlowStatusOk(cloneCtx, flow.OrderID, flow.Uid, int64(st.OrderStatusExpireDay)*86400)
		if err != nil {
			logger.Error(cloneCtx, "afterScoreOp dao.MarkOrderFlowStatusOk fail", zap.Error(err))
			// 这里不影响主进程
		}
		return nil
	}, nil)

	return nil
}

type ScoreChangeSideEffect struct {
	side_effect.BaseSideEffect
}

func (ScoreChangeSideEffect) ScoreChange(ctx context.Context, st *model.ScoreType, flow *dao.ScoreFlowModel) error {
	if !conf.Conf.WriteScoreFlow {
		return nil
	}

	// 写入流水
	err := writeScoreFlow(ctx, st, flow)
	if err != nil {
		logger.Error(ctx, "SideEffect.ScoreChange call writeScoreFlow fail", zap.Any("flow", flow), zap.Error(err))
	}
	return nil
}
