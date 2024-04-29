package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func InitializeRoutes(Fiber *fiber.App) {

	Fiber.Use(cors.New())
	Fiber.Use(etag.New())
	Fiber.Use(logger.New(logger.Config{
		// For more options, see the Config section
		Format:     "${pid} ${time} ${locals:requestid} ${status} - ${method} ${path}\n",
		TimeFormat: "02-Jan-2006",
		TimeZone:   "America/Sao_Paulo",
	}))
	Fiber.Get("/metrics", monitor.New())

	Fiber.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	api := Fiber.Group("/api")

	api.Get("/info", getInfo)
	api.Get("/info/ps", getInfoProcess)
	api.Get("/gpio", getGpio)
	api.Patch("/gpio/:pin", updateGpio)

	api.Get("/share", getShare)
	api.Get("/share/*", getShareFile)

	//auth := api.Group("/auth")

}
