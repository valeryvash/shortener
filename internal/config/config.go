package config

import (
	"flag"
	"fmt"
)

type Config struct {
	BindAddress     string
	RedirectAddress string
}

func NewConfig(args []string) (Config, error) {
	var config Config

	set := flag.NewFlagSet("shortener", flag.ContinueOnError)

	set.StringVar(&config.BindAddress, "a", "localhost:8080", "Bind address for server start")
	set.StringVar(&config.RedirectAddress, "b", "http://localhost:8080", "Redirect base address")

	err := set.Parse(args)
	if err != nil {
		return config, fmt.Errorf("failed to parse args: %w", err)
	}

	return config, nil
}
