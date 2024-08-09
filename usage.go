package score

import (
	"github.com/zly-app/zapp"
	"github.com/zly-app/zapp/core"
)

func WithService() zapp.Option {
	return zapp.WithCustomEnableService(func(app core.IApp, services []core.ServiceType) []core.ServiceType {
		return services
	})
}
