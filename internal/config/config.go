package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Mongo    Mongo
	AWS      AWS
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

type AWS struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	BucketName      string
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
		return err
	}

	cfg.Mongo.URI = os.Getenv("MONGO_INITDB_ROOT_URI")
	cfg.Mongo.Database = os.Getenv("MONGO_INITDB_ROOT_DATABASE")
	cfg.Mongo.Username = os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	cfg.Mongo.Password = os.Getenv("MONGO_INITDB_ROOT_PASSWORD")

	cfg.AWS.AccessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")
	cfg.AWS.SecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	cfg.AWS.Region = os.Getenv("AWS_DEFAULT_REGION")
	cfg.AWS.BucketName = os.Getenv("AWS_S3_BUCKET_NAME")

	cfg.Password.Salt = os.Getenv("PASSWORD_SALT")

	cfg.JWT.Secret = os.Getenv("JWT_SECRET")

	return nil
}
