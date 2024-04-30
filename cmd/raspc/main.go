package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gabrielmoura/raspController/configs"
	"github.com/gabrielmoura/raspController/infra/db"
	"github.com/gabrielmoura/raspController/infra/gpio"
	"github.com/gabrielmoura/raspController/infra/routes"
	"github.com/gabrielmoura/raspController/pkg/mdns"
	"github.com/gofiber/fiber/v2"
)

func main() {
	ctx := context.Background()

	if err := configs.LoadConfig(); err != nil {
		panic(err)
	}

	app := fiber.New(
		fiber.Config{
			Prefork:       false,
			CaseSensitive: true,
			StrictRouting: true,
			ServerHeader:  configs.Conf.AppName,
			AppName:       configs.Conf.AppName,
		},
	)
	routes.InitializeRoutes(app)
	db.Initialize(ctx)

	gpio.Initialize(ctx)
	if err := mdns.SetDNS(configs.Conf.AppName, configs.Conf.Port); err != nil {
		log.Println("Error setting DNS:", err)
	}

	log.Fatal(app.Listen(fmt.Sprintf(":%d", configs.Conf.Port)))
}
