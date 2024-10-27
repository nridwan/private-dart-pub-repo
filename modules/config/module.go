package config

import (
	"log"
	"private-pub-repo/base"

	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

type ConfigModule struct {
}

func NewModule() *ConfigModule {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return &ConfigModule{}
}

func ProvideService(module *ConfigModule) ConfigService {
	return module
}

func fxRegister(lifeCycle fx.Lifecycle, module *ConfigModule) {
	base.FxRegister(module, lifeCycle)
}

func SetupModule() *ConfigModule {
	return NewModule()
}

var FxModule = fx.Module("Config", fx.Provide(NewModule), fx.Provide(ProvideService), fx.Invoke(fxRegister))

// implements `BaseModule` of `base/module.go` start

func (module *ConfigModule) OnStart() error {
	return nil
}

func (module *ConfigModule) OnStop() error {
	return nil
}

// implements `BaseModule` of `base/module.go` end
