package monitor

import (
	"private-pub-repo/base"
	"private-pub-repo/modules/app"
	"private-pub-repo/modules/config"

	"github.com/gofiber/fiber/v2"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"
)

type MonitorModule struct {
	Service MonitorService
	app     *fiber.App
	config  config.ConfigService
	tp      *sdktrace.TracerProvider
}

func NewModule(service MonitorService, app *fiber.App, config config.ConfigService) *MonitorModule {
	return &MonitorModule{Service: service, app: app, config: config}
}

func fxRegister(lifeCycle fx.Lifecycle, module *MonitorModule) {
	base.FxRegister(module, lifeCycle)
}

func SetupModule(app *app.AppModule, config *config.ConfigModule) *MonitorModule {
	return NewModule(provideMonitorService(), app.App, config)
}

var FxModule = fx.Module("Monitor", fx.Provide(NewModule), fx.Provide(provideMonitorService), fx.Invoke(fxRegister))

// implements `BaseModule` of `base/module.go` start

func (module *MonitorModule) OnStart() error {
	module.initOpentelemetry()
	module.registerRoutes()
	return nil
}

func (module *MonitorModule) OnStop() error {
	module.destroyOpentelemetry()
	return nil
}

// implements `BaseModule` of `base/module.go` end
