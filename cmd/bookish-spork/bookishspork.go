package main

import (
	"fmt"
	"os"

	"github.com/dstdfx/bookish-spork/cmd/bookish-spork/app"
)

func main() {
	if err := app.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
