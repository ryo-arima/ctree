package main

import (
	"log"

	"github.com/ryo-arima/ctree/pkg"
	"github.com/ryo-arima/ctree/pkg/config"
)

func main() {
	// Load configuration
	conf := config.NewConfig()

	// Try to load config from file if exists
	configPath := config.GetConfigPath()
	if configPath != "" {
		loadedConf, err := config.LoadConfigFromFile(configPath)
		if err != nil {
			log.Printf("Warning: Failed to load config from %s: %v", configPath, err)
		} else {
			conf = loadedConf
		}
	}

	// Start the CLI
	pkg.ClientForCtree(conf)
}
