package config

import (
	"fmt"
	"nevermore/internal/storage/minio"
	"nevermore/pkg/logger"

	"nevermore/internal/storage/postgres"
	"nevermore/internal/storage/redis"

	"github.com/spf13/viper"
)

func Psql() postgres.Config {
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

func Rds() redis.Config {
	result := redis.Config{
		User:     viper.GetString("redis.user"),
		Password: viper.GetString("redis.password"),
		Url:      viper.GetString("redis.url"),
		DB:       viper.GetInt("redis.db"),
	}

	return result
}

func JwtSecret() string {
	return viper.GetString("jwt.secret")
}

func Photoes() minio.Config {
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

func Srv() string {
	return fmt.Sprintf(":%d", viper.GetInt("server.port"))
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
