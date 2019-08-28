package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"lll.github.com/llleaas/cmd/hermes/app"
	"lll.github.com/llleaas/pkg/hermes/server"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	command := app.NewHermesCommand(server.SetupSignalHandler())
	if err := command.Execute(); err != nil {
		fmt.Fprint(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
