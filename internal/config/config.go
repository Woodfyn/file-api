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
	Firebase Firebase
	Server   Server
	JWT      JWT
	Password Password
}

type Mongo struct {
	URI      string
	Database string
	Username string
	Password string
}

type RDB struct {
	Addr string
}

type Firebase struct {
	FileName   string
	BucketName string
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

	cfg.Mongo.URI = os.Getenv("MONGO_INITDB_ROOT_URI")
	cfg.Mongo.Database = os.Getenv("MONGO_INITDB_ROOT_DATABASE")
	cfg.Mongo.Username = os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	cfg.Mongo.Password = os.Getenv("MONGO_INITDB_ROOT_PASSWORD")

	cfg.RDB.Addr = os.Getenv("REDIS_PORT")

	cfg.Firebase.FileName = os.Getenv("FIREBASE_FILE_NAME")
	cfg.Firebase.BucketName = os.Getenv("FIREBASE_BUCKET_NAME")

	cfg.Password.Salt = os.Getenv("PASSWORD_SALT")

	cfg.JWT.Secret = os.Getenv("JWT_SECRET")

	return nil
}
