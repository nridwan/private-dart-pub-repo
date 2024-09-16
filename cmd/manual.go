package cmd

import (
	"gofiber-boilerplate/base"
	"gofiber-boilerplate/modules/app"
	"gofiber-boilerplate/modules/config"
	"gofiber-boilerplate/modules/db"
	"gofiber-boilerplate/modules/jwt"
	"gofiber-boilerplate/modules/monitor"
	"gofiber-boilerplate/modules/user"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/urfave/cli/v2"
)

func CommandManual() *cli.Command {
	return &cli.Command{
		Name:  "manual",
		Usage: "start manual server",
		Action: func(cCtx *cli.Context) error {
			runManual()
			return nil
		},
	}
}

func runManual() {
	configModule := config.SetupModule()
	appModule := app.SetupModule(configModule)
	monitorModule := monitor.SetupModule(appModule, configModule)
	dbModule := db.SetupModule(configModule)
	jwtModule := jwt.SetupModule(appModule, configModule)
	userModule := user.SetupModule(appModule, dbModule, jwtModule, monitorModule)

	modules := []base.BaseModule{
		configModule,
		appModule,
		monitorModule,
		dbModule,
		jwtModule,
		userModule,
	}

	for i := range modules {
		modules[i].OnStart()
	}

	appModule.App.Get("/", func(c *fiber.Ctx) error {
		return fiber.NewError(400, "Error")
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() error {
		<-c
		log.Println("stopping server ...")
		go func() {
			appModule.App.Shutdown()
		}()
		for i := range modules {
			modules[i].OnStop()
		}
		log.Println("stop server success")
		return nil
	}()

	// ...

	if err := appModule.App.Listen(configModule.Getenv("APP_HOST", "") + ":" + configModule.Getenv("PORT", "3000")); err != nil {
		log.Panic(err)
	}
}
