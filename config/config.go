package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App             `yaml:"app"`
		HTTP            `yaml:"http"`
		Log             `yaml:"logger"`
		Cache           `yaml:"cache"`
		Repository      `yaml:"repository"`
		DiskRepository  `yaml:"diskRepository"`
		MongoRepository `yaml:"mongoRepository"`
	}

	// App -.
	App struct {
		Name    string `yaml:"name" env:"APP_NAME" validate:"required"`
		Version string `yaml:"version" env:"APP_VERSION" validate:"required"`
	}

	// HTTP -.
	HTTP struct {
		Port string `yaml:"port" env:"HTTP_PORT" validate:"required"`
	}

	// Log -.
	Log struct {
		Level string `yaml:"log_level" env:"LOG_LEVEL" validate:"required"`
	}

	Cache struct {
		Size int `yaml:"size" env:"CACHE_SIZE" validate:"required"`
	}

	Repository struct {
		Type string `yaml:"type" env:"REPOSITORY_TYPE" validate:"required,oneof=disk mongo"`
	}

	DiskRepository struct {
		RelativePath string `yaml:"relativePath" env:"DISK_REPOSITORY_RELATIVE_PATH"`
	}

	MongoRepository struct {
		URI        string `yaml:"uri" env:"MONGO_REPOSITORY_URI"`
		DB         string `yaml:"db" env:"MONGO_REPOSITORY_DB"`
		Collection string `yaml:"collection" env:"MONGO_REPOSITORY_COLLECTION"`
	}
)

// NewConfig returns app config.
func NewConfig(path string) (*Config, error) {
	cfg := &Config{}

	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, err
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	return cfg, nil
}
