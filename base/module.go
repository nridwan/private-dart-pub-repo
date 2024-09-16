package base

import (
	"context"

	"go.uber.org/fx"
)

type BaseModule interface {
	OnStart() error
	OnStop() error
}

func FxRegister(module BaseModule, lifeCycle fx.Lifecycle) {
	lifeCycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			return module.OnStart()
		},
		OnStop: func(_ context.Context) error {
			return module.OnStop()
		},
	})
}
