package app

import (
	"fmt"

	"github.com/quinncuatro/hass-cli/internal/cli"
	"github.com/quinncuatro/hass-cli/internal/config"
)

func Run(args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	commander := cli.NewCommander(cfg)
	return commander.Execute(args)
}