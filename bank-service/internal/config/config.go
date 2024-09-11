package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		GRPC     GRPCConfig     `yaml:"grpc"`
		Postgres PostgresConfig `yaml:"postgres"`
	}

	GRPCConfig struct {
		Port    int           `env-required:"true" yaml:"port" env:"GRPC_PORT"`
		Timeout time.Duration `env-required:"true" yaml:"timeout" env:"GRPC_TIMEOUT"`
	}

	PostgresConfig struct {
		URL         string `env-required:"true" yaml:"url" env:"PG_URL"`
		MaxPoolSize int    `yaml:"max_pool_size" env:"PG_MAX_POOL_SIZE"`
	}
)

func MustLoad() *Config {
	var cfg Config
	path := fetchConfigPath()
	if path != "" {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			panic("config file doesn't exist: " + path)
		}
		if err := cleanenv.ReadConfig(path, &cfg); err != nil {
			panic("failed to read config file: " + err.Error())
		}
		return &cfg
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("failed to read configuration from env")
	}

	return &cfg
}

func fetchConfigPath() string {
	var result string
	flag.StringVar(&result, "config", "", "path to config file")
	flag.Parse()
	if result == "" {
		result = os.Getenv("CONFIG_PATH")
	}
	return result
}
