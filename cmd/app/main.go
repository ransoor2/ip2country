package main

import (
	"log"

	"github.com/ransoor2/ip2country/config"
	"github.com/ransoor2/ip2country/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig("./config/config.yml")
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
