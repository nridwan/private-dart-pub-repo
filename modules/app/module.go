package app

import (
	"private-pub-repo/base"
	"private-pub-repo/modules/config"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

type AppModule struct {
	App             *fiber.App
	ResponseService ResponseService
	Validator       *validator.Validate
}

func NewFiber(responseService ResponseService) *fiber.App {
	return fiber.New(fiber.Config{
		ErrorHandler: responseService.ErrorHandler,
	})
}

func ProvideValidator() *validator.Validate {
	return validator.New()
}

func NewModule(app *fiber.App, responseService ResponseService, validator *validator.Validate) *AppModule {
	return &AppModule{
		App:             app,
		ResponseService: responseService,
		Validator:       validator,
	}
}

func fxRegister(lifeCycle fx.Lifecycle, module *AppModule) {
	base.FxRegister(module, lifeCycle)
}

func SetupModule(config config.ConfigService) *AppModule {
	responseService := NewResponseService(config)
	return NewModule(NewFiber(responseService), responseService, ProvideValidator())
}

var FxModule = fx.Module("app", fx.Provide(NewFiber), fx.Provide(NewResponseService), fx.Provide(ProvideValidator), fx.Provide(NewModule), fx.Invoke(fxRegister))

// implements `BaseModule` of `base/module.go` start

func (module *AppModule) OnStart() error {
	return nil
}

func (module *AppModule) OnStop() error {
	return nil
}

// implements `BaseModule` of `base/module.go` end
