package main

import (
	"fmt"
	"github.com/gabrielmoura/raspController/configs"
	"github.com/gabrielmoura/raspController/infra/gpio"
	"github.com/gabrielmoura/raspController/infra/routes"
	"github.com/gabrielmoura/raspController/pkg/mdns"
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {
	app := fiber.New()
	routes.InitializeRoutes(app)

	if err := configs.LoadConfig(); err != nil {
		panic(err)
	}
	gpio.Initialize()
	if err := mdns.SetDNS(configs.Conf.AppName, configs.Conf.Port); err != nil {
		log.Println("Error setting DNS:", err)
	}

	log.Fatal(app.Listen(fmt.Sprintf(":%d", configs.Conf.Port)))
}
