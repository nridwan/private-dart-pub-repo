package cmd

import (
	"log"
	"os"
	"os/signal"
	"private-pub-repo/base"
	"private-pub-repo/modules/app"
	"private-pub-repo/modules/config"
	"private-pub-repo/modules/db"
	"private-pub-repo/modules/jwt"
	"private-pub-repo/modules/monitor"
	"private-pub-repo/modules/user"
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