package main

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"log"

	"github.com/ttmbank/backend/adminserver/app"
)

// main
func main() {
	ctx := context.Background()
	viper.SetConfigFile("./config/config.yml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Config error", err)
	}

	port := viper.GetUint("port.adminapi")
	bindHost := fmt.Sprintf(":%d", port)

	app := &app.App{}
	app.Initialize(ctx)
	app.Run(ctx, bindHost)
}
