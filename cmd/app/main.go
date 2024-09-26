package main

import (
	"flag"
	"log"

	"github.com/ransoor2/ip2country/config"
	"github.com/ransoor2/ip2country/internal/app"
)

func main() {
	// Define a flag for the relative path
	relativePath := flag.String("config-path", "./config/config.yml", "Path to the configuration file")

	// Parse the flags
	flag.Parse()

	// Configuration
	cfg, err := config.NewConfig(*relativePath)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
