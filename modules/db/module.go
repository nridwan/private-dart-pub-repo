package db

import (
	"gofiber-boilerplate/base"
	"gofiber-boilerplate/modules/config"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

type DbModule struct {
	config      config.ConfigService
	db          map[string]*gorm.DB
	autoMigrate bool
}

func NewModule(config config.ConfigService) *DbModule {
	return &DbModule{config: config, db: map[string]*gorm.DB{}}
}

func ProvideService(module *DbModule) DbService {
	return module
}

func fxRegister(lifeCycle fx.Lifecycle, module *DbModule) {
	base.FxRegister(module, lifeCycle)
}

func SetupModule(config config.ConfigService) *DbModule {
	return NewModule(config)
}

var FxModule = fx.Module("Db", fx.Provide(NewModule), fx.Provide(ProvideService), fx.Invoke(fxRegister))

// implements `BaseModule` of `base/module.go` start

func (module *DbModule) OnStart() error {
	module.addDefaultConfig()
	return nil
}

func (module *DbModule) OnStop() error {
	return nil
}

// implements `BaseModule` of `base/module.go` end
