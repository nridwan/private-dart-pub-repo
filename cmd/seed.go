package cmd

import (
	"context"
	"os"
	"private-pub-repo/modules/app"
	"private-pub-repo/modules/config"
	"private-pub-repo/modules/db"
	"private-pub-repo/modules/jwt"
	"private-pub-repo/modules/monitor"
	"private-pub-repo/modules/user"

	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
)

func CommandDbSeed() *cli.Command {
	return &cli.Command{
		Name:  "db:seed",
		Usage: "apply seeder",
		Action: func(cCtx *cli.Context) error {
			runSeeder()
			return nil
		},
	}
}

func runSeeder() {
	fxApp := fx.New(
		config.FxModule,
		app.FxModule,
		monitor.FxModule,
		db.FxModule,
		jwt.FxModule,
		user.FxModule,
		fx.Invoke(applySeeders),
		fx.NopLogger,
	)

	fxApp.Run()
}

func applySeeders(
	lifeCycle fx.Lifecycle,
	userModule *user.UserModule,
) {
	lifeCycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			userModule.RunSeeder()
			os.Exit(0)
			return nil
		},
		OnStop: func(_ context.Context) error {
			return nil
		},
	})
}
