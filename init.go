package score

import (
	"github.com/zly-app/zapp"
	"github.com/zly-app/zapp/config"
	"github.com/zly-app/zapp/core"
	"github.com/zly-app/zapp/handler"
	"go.uber.org/zap"

	"github.com/zlyuancn/score/conf"
	"github.com/zlyuancn/score/dao"
	"github.com/zlyuancn/score/score_flow"
	"github.com/zlyuancn/score/score_type"
	"github.com/zlyuancn/score/side_effect"
)

func init() {
	config.RegistryApolloNeedParseNamespace(conf.ScoreConfigKey)

	// 注册副作用
	side_effect.RegistrySideEffect("score_flow", new(score_flow.ScoreChangeSideEffect))

	// 持久内存-加载积分类型
	score_type.StartLoopLoad()

	zapp.AddHandler(zapp.BeforeInitializeHandler, func(app core.IApp, handlerType handler.HandlerType) {
		err := app.GetConfig().Parse(conf.ScoreConfigKey, &conf.Conf, true)
		if err != nil {
			app.Fatal("parse score config err", zap.Error(err))
		}
		conf.Conf.Check()
	})
	zapp.AddHandler(zapp.AfterInitializeHandler, func(app core.IApp, handlerType handler.HandlerType) {
		dao.TryInjectScript()
	})
}
