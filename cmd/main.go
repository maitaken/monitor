package main

import (
	"fmt"
	"os"

	"github.com/maitaken/monitor/app"
	"github.com/maitaken/monitor/app/config"
)

func main() {
	c, err := config.NewConfig()
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	if err := app.Run(c); err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
}
