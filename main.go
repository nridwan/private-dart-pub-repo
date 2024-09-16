package main

import (
	"fmt"
	"log"
	"os"
	"private-pub-repo/cmd"
	"runtime"

	"github.com/urfave/cli/v2"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	commands := []*cli.Command{
		cmd.CommandManual(),
		cmd.CommandFx(),
		cmd.CommandDbSeed(),
	}

	app := &cli.App{
		Commands: commands,
		Name:     "apiserver",
		Usage:    "manual, fx, db:seed",
		Action: func(cli *cli.Context) error {
			fmt.Printf("%s version:%s\n", cli.App.Name, "3.0")
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
