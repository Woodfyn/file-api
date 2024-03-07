package config

import (
	"os"
	"time"

	"github.com/Woodfyn/file-api/internal/core"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Mongo    Mongo
	RDB      RDB
	Server   Server
	JWT      JWT
	Password Password
}

type Mongo struct {
	URI      string
	Username string
	Password string
}

type RDB struct {
	ADR      string
	Password string
}

type Server struct {
	Port string
}

type JWT struct {
	AccessTokenTTL  time.Duration `mapstructure:"accsess_token_ttl"`
	RefreshTokenTTL time.Duration `mapstructure:"refresh_token_ttl"`
	Secret          string
}

type Password struct {
	Salt string
}

func InitConfig(folder, file string) (*Config, error) {
	cfg := new(Config)

	viper.AddConfigPath(folder)
	viper.SetConfigName(file)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	err := setFromEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func setFromEnv(cfg *Config) error {
	err := godotenv.Load()
	if err != nil {
		return core.ErrFileNotFound
	}

	cfg.Mongo.URI = os.Getenv("MONGO_URI")
	cfg.Mongo.Username = os.Getenv("MONGO_USERNAME")
	cfg.Mongo.Password = os.Getenv("MONGO_PASSWORD")

	cfg.RDB.ADR = os.Getenv("REDIS_ADR")
	cfg.RDB.Password = os.Getenv("REDIS_PASSWORD")

	cfg.Password.Salt = os.Getenv("PASSWORD_SALT")

	cfg.JWT.Secret = os.Getenv("JWT_SECRET")

	return nil
}
