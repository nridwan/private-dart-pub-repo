package jwt

import (
	"gofiber-boilerplate/base"
	"gofiber-boilerplate/modules/app"
	"gofiber-boilerplate/modules/config"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/fx"
)

type JwtClaim = jwt.MapClaims

type JwtModule struct {
	config          config.ConfigService
	responseService app.ResponseService
	lifetime        time.Duration
	refreshLifetime time.Duration
	secret          string
	handler         fiber.Handler
}

func NewModule(config config.ConfigService, responseService app.ResponseService) *JwtModule {
	return &JwtModule{
		config:          config,
		responseService: responseService,
	}
}

func ProvideService(module *JwtModule) JwtService {
	return module
}

func fxRegister(lifeCycle fx.Lifecycle, module *JwtModule) {
	base.FxRegister(module, lifeCycle)
}

func SetupModule(app *app.AppModule, config *config.ConfigModule) *JwtModule {
	return NewModule(config, app.ResponseService)
}

var FxModule = fx.Module("Jwt", fx.Provide(NewModule), fx.Provide(ProvideService), fx.Invoke(fxRegister))

// implements `BaseModule` of `base/module.go` start

func (module *JwtModule) OnStart() error {
	module.Init(module.config)

	return nil
}

func (module *JwtModule) OnStop() error {
	return nil
}

// implements `BaseModule` of `base/module.go` end
