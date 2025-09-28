package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"syscall"

	_ "nevermore/docs"
	"nevermore/internal/app"
	exit "nevermore/pkg/context"

	"github.com/spf13/viper"
)

//export GOOSE_DBSTRING="user=postgres password=1 dbname=nevermore sslmode=disable
//echo $GOOSE_DBSTRING

//export GOOSE_DRIVER=postgres
//echo $GOOSE_DRIVER

//sudo -u postgres psql
//\c networks

//swag init --generalInfo cmd/Cringe-Networks/main.go --output docs

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	fmt.Println("Loading config...")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Ошибка чтения конфигурации: %s", err)
	}
}

// @title		Nevermore API
// @version		1.0
// @description	API для Nevermore
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
