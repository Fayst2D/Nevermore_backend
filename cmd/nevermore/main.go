package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"nevermore/config"
	"nevermore/pkg/logger"
	"os"
	"strings"
	"syscall"

	_ "nevermore/docs"
	"nevermore/internal/app"
	exit "nevermore/pkg/context"

	"github.com/spf13/viper"
)

func init() {
	_ = godotenv.Load()

	viper.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))
	viper.AutomaticEnv()

	if err := logger.Init(config.NewLogger()); err != nil {
		panic(err)
	}
}

// @title		Backend API
// @version		1.0
// @description	API для backend
// @host		localhost:3000
// @BasePath	/
// @schemes		http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введите токен в формате: Bearer {token}

// @tag.name user
// @tag.description Операции с пользователями
func main() {
	app, err := app.New()
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx, cancel := exit.WithSignal(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := app.Run(ctx); err != nil {
		fmt.Println(err)
		return
	}

	return
}
