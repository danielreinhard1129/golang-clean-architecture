package main

import (
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/danielreinhard1129/fiber-clean-arch/internal/entities"
)

func main() {

	stmts, err := gormschema.New("postgres").Load(&entities.User{}) // add new entities here
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)
}
