package config

import (
	"fmt"
	"log"
	"nevermore/pkg/logger"

	"github.com/joho/godotenv"

	"nevermore/internal/storage/postgres"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port int    `mapstructure:"port"`
		Host string `mapstructure:"host"`
	} `mapstructure:"server"`
	Postgres struct {
		Url      string `mapstructure:"url"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Name     string `mapstructure:"name"`
		Driver   string `mapstructure:"driver"`
	} `mapstructure:"postgres"`
	Redis struct {
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Url      string `mapstructure:"url"`
		DB       int    `mapstructure:"db"`
	} `mapstructure:"redis"`
	Minio struct {
		Endpoint  string `mapstructure:"endpoint"`
		AccessKey string `mapstructure:"access_key"`
		SecretKey string `mapstructure:"secret_key"`
		Photoes   string `mapstructure:"photoes"`
		Pages     string `mapstructure:"pages"`
		Pdfs      string `mapstructure:"pdfs"`
	} `mapstructure:"minio"`
	Logger struct {
		Dir               string `mapstructure:"dir"`
		Filename          string `mapstructure:"filename"`
		Level             string `mapstructure:"level"`
		MaxSizeMB         int    `mapstructure:"max_size_mb"`
		MaxBackups        int    `mapstructure:"max_backups"`
		MaxAgeDays        int    `mapstructure:"max_age_days"`
		Compress          bool   `mapstructure:"compress"`
		DuplicateToStdout bool   `mapstructure:"duplicate_to_stdout"`
		TimeFormat        string `mapstructure:"time_format"`
		ServiceName       string `mapstructure:"service_name"`
	} `mapstructure:"logger"`
}

func (c Config) Psql() postgres.Config {
	result := postgres.Config{
		URL: fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable",
			c.Postgres.User,
			c.Postgres.Password,
			c.Postgres.Url,
			c.Postgres.Name,
		),

		Driver: c.Postgres.Driver,
	}

	return result
}

func (c Config) Srv() string {
	return fmt.Sprintf(":%d", c.Server.Port)
}

func Init() (Config, error) {
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Ошибка декодирования конфигурации: %s", err)
	}

	return config, nil
}

func (c Config) NewLogger() logger.Config {
	return logger.Config{
		Dir:               c.Logger.Dir,
		Filename:          c.Logger.Filename,
		Level:             c.Logger.Level,
		MaxSizeMB:         c.Logger.MaxSizeMB,
		MaxBackups:        c.Logger.MaxBackups,
		MaxAgeDays:        c.Logger.MaxAgeDays,
		Compress:          c.Logger.Compress,
		DuplicateToStdout: c.Logger.DuplicateToStdout,
		TimeFormat:        c.Logger.TimeFormat,
		ServiceName:       c.Logger.ServiceName,
	}
}

func LoadConfig() (Config, error) {
	_ = godotenv.Load()

	viper.AutomaticEnv()

	viper.BindEnv("app.name", "APP_NAME")
	viper.BindEnv("app.version", "APP_VERSION")
	viper.BindEnv("app.port", "APP_PORT")

	viper.BindEnv("database.url", "DB_HOST")
	viper.BindEnv("database.user", "DB_USER")
	viper.BindEnv("database.password", "DB_PASSWORD")
	viper.BindEnv("database.name", "DB_NAME")
	viper.BindEnv("database.driver", "DB_DRIVER")

	viper.BindEnv("minio.endpoint", "MINIO_ENDPOINT")
	viper.BindEnv("minio.access_key", "MINIO_ACCESS_KEY")
	viper.BindEnv("minio.secret_key", "MINIO_SECRET_KEY")
	viper.BindEnv("minio.bucket", "MINIO_BUCKET")

	viper.BindEnv("redis.url", "REDIS_HOST")
	viper.BindEnv("redis.user", "REDIS_USER")
	viper.BindEnv("redis.password", "REDIS_PASSWORD")
	viper.BindEnv("redis.db", "REDIS_DB")

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Ошибка декодирования конфигурации: %s", err)
	}

	return cfg, nil
}
