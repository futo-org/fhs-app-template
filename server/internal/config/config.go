package config

import (
	"fmt"
	"os"

	"go.uber.org/fx"
)

type Config struct {
	DatabasePath string
	AppEnv       string
	ServerAddr   string
	WebRoot      string
}

func NewConfig() (*Config, error) {
	appEnv := os.Getenv("APP_ENV")
	fmt.Printf("environment is %s\n", appEnv)

	serverAddr := os.Getenv("SERVER_ADDR")
	if serverAddr == "" {
		serverAddr = "0.0.0.0:3001"
	}

	webRoot := os.Getenv("WEB_ROOT")
	if webRoot == "" {
		webRoot = "web"
	}

	return &Config{
		DatabasePath: "appdata/database.db",
		AppEnv:       appEnv,
		ServerAddr:   serverAddr,
		WebRoot:      webRoot,
	}, nil
}

var Module = fx.Module("config", fx.Provide(NewConfig))
