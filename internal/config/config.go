package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/fx"
)

const (
	defaultPath = "./config/config.yml"

	environmentVariable = "HARMONY_CONFIG"
)

type Config struct {
	fx.Out

	*App        `yaml:"app"`
	*Logger     `yaml:"logger"`
	*Http       `yaml:"http"`
	*Jwt        `yaml:"jwt"`
	*Mongo      `yaml:"mongo"`
	*Centrifugo `yaml:"centrifugo"`
}

type App struct {
	Development bool `yaml:"development"`
}

type Logger struct {
	Level string `yaml:"level"`
}

type Jwt struct {
	Issuer         string        `yaml:"issuer"`
	Lifetime       time.Duration `yaml:"lifetime"`
	PrivateKeyPath string        `yaml:"private_key_path"`
}

type Http struct {
	Address      string        `yaml:"address"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`

	// CORS settings
	CorsAllowOrigins     []string `yaml:"cors_allow_origins"`
	CorsAllowCredentials bool     `yaml:"cors_allow_credentials"`
}

type Mongo struct {
	Address  string `yaml:"address"`
	Direct   bool   `yaml:"direct"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type Centrifugo struct {
	ApiAddress string `yaml:"api_address"`
	ApiKey     string `yaml:"api_key"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := cleanenv.ReadConfig(getConfigPath(), &config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func getConfigPath() string {
	envPath := os.Getenv(environmentVariable)
	if envPath != "" {
		return envPath
	}

	flagPath := flag.String("config", defaultPath, "path to configuration file")

	// defaultPath will be returned if no flag provided
	return *flagPath
}
