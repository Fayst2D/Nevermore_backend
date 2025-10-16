package config

import (
	"fmt"
	"log"
	"nevermore/internal/storage/minio"
	"nevermore/pkg/logger"

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
	Minio struct {
		Endpoint  string `mapstructure:"endpoint"`
		AccessKey string `mapstructure:"access_key"`
		SecretKey string `mapstructure:"secret_key"`
		Photoes   string `mapstructure:"photoes"`
		Pages     string `mapstructure:"pages"`
		Pdfs      string `mapstructure:"pdfs"`
	} `mapstructure:"minio"`
}

func (c Config) Psql() postgres.Config {
	result := postgres.Config{
		URL: fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable",
			viper.GetString("postgres.user"),
			viper.GetString("postgres.password"),
			viper.GetString("postgres.url"),
			viper.GetString("postgres.name"),
		),
		Driver: viper.GetString("postgres.driver"),
	}

	return result
}

func (c Config) Photoes() minio.Config {
	result := minio.Config{
		AccessKeyID:     viper.GetString("minio.access_key"),
		SecretAccessKey: viper.GetString("minio.secret_key"),
		BaseURL:         viper.GetString("minio.endpoint"),
		Photoes:         viper.GetString("minio.photoes"),
		Pdfs:            viper.GetString("minio.pdfs"),
		Pages:           viper.GetString("minio.pages"),
	}

	return result
}

func (c Config) Srv() string {
	return fmt.Sprintf(":%d", viper.GetInt("server.port"))
}

func Init() (Config, error) {
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Ошибка декодирования конфигурации: %s", err)
	}

	return config, nil
}

func NewLogger() logger.Config {
	return logger.Config{
		Dir:               viper.GetString("logger.dir"),
		Filename:          viper.GetString("logger.filename"),
		Level:             viper.GetString("logger.level"),
		MaxSizeMB:         viper.GetInt("logger.max_size_mb"),
		MaxBackups:        viper.GetInt("logger.max_backups"),
		MaxAgeDays:        viper.GetInt("logger.max_age_days"),
		Compress:          viper.GetBool("logger.compress"),
		DuplicateToStdout: viper.GetBool("logger.duplicate_to_stdout"),
		TimeFormat:        viper.GetString("logger.time_format"),
		ServiceName:       viper.GetString("logger.service_name"),
	}
}
