package score

import (
	"github.com/zly-app/zapp"
	"github.com/zly-app/zapp/config"
	"github.com/zly-app/zapp/core"
	"github.com/zly-app/zapp/handler"
	"go.uber.org/zap"

	"github.com/zlyuancn/score/conf"
	"github.com/zlyuancn/score/score_type"
)

func init() {
	config.RegistryApolloNeedParseNamespace(conf.ScoreConfigKey)

	score_type.Init()

	zapp.AddHandler(zapp.BeforeInitializeHandler, func(app core.IApp, handlerType handler.HandlerType) {
		err := app.GetConfig().Parse(conf.ScoreConfigKey, &conf.Conf, true)
		if err != nil {
			app.Fatal("parse score config err", zap.Error(err))
		}
		conf.Conf.Check()
	})
}
