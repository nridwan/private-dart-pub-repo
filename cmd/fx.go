package cmd

import (
	"context"
	"log"
	"private-pub-repo/modules/app"
	"private-pub-repo/modules/config"
	"private-pub-repo/modules/db"
	"private-pub-repo/modules/jwt"
	"private-pub-repo/modules/monitor"
	"private-pub-repo/modules/pub"
	"private-pub-repo/modules/pubtoken"
	"private-pub-repo/modules/user"

	"github.com/gofiber/fiber/v2"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
)

func CommandFx() *cli.Command {
	return &cli.Command{
		Name:  "fx",
		Usage: "start fx server",
		Action: func(cCtx *cli.Context) error {
			runFx()
			return nil
		},
	}
}

func runFx() {
	fxApp := fx.New(
		config.FxModule,
		app.FxModule,
		monitor.FxModule,
		db.FxModule,
		jwt.FxModule,
		user.FxModule,
		pubtoken.FxModule,
		pub.FxModule,
		fx.Invoke(registerWebServer),
	)

	fxApp.Run()
}

func registerWebServer(
	lifeCycle fx.Lifecycle,
	app *fiber.App,
	config config.ConfigService,
) {
	lifeCycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				app.Get("/", func(c *fiber.Ctx) error {
					return fiber.NewError(400, "Error")
				})

				if err := app.Listen(
					config.Getenv("APP_HOST", "") + ":" + config.Getenv("PORT", "3000"),
				); err != nil {
					log.Fatalf("start server error : %v\n", err)
				}

			}()
			return nil
		},
		OnStop: func(_ context.Context) error {
			log.Println("stopping server ...")
			err := app.Shutdown()
			log.Println("stop server success")
			return err
		},
	})
}
