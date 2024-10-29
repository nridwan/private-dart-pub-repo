package main

import (
	"fmt"
	"io"
	"os"
	"private-pub-repo/modules/pub/pubmodel"
	"private-pub-repo/modules/pubtoken/pubtokenmodel"
	"private-pub-repo/modules/user/usermodel"

	"ariga.io/atlas-provider-gorm/gormschema"
)

func main() {
	stmts, err := gormschema.New("postgres").Load(
		// user module
		&usermodel.UserModel{},
		// pubtoken module
		&pubtokenmodel.PubTokenModel{},
		// pub module
		&pubmodel.PubPackageModel{},
		&pubmodel.PubVersionModel{},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	io.WriteString(os.Stdout, stmts)
}
