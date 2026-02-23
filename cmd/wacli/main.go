package main

import (
	"os"

	"github.com/steipete/wacli/internal/config"
)

func main() {
	config.ApplyDeviceConfigFromEnv()
	if err := execute(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		os.Exit(1)
	}
}
