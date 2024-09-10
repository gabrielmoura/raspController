package routes

import (
	"github.com/gabrielmoura/raspController/configs"
	"github.com/gabrielmoura/raspController/infra/middleware"
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
		TimeFormat: configs.Conf.TimeFormat,
		TimeZone:   configs.Conf.TimeZone,
	}))
	Fiber.Get("/metrics", monitor.New())

	Fiber.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("RaspController API")
	})
	api := Fiber.Group("/api")

	api.Get("", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"/api/info":      "Returns system information.",
			"/api/info/ps":   "Returns process information.",
			"/api/info/net":  "Returns network information.",
			"/api/info/mem":  "Returns memory information.",
			"/api/info/disk": "Returns disk information.",
			"/api/info/gpio": "Returns list of available GPIOs",
			"/api/info/usb":  "Returns list of USB devices",
			"/api/info/cpu":  "Returns CPU information.",
			"/api/gpio":      "Returns the status of all configured GPIO pins.",
			"/api/gpio/all":  "Returns all GPIO pins from the GPIO chip.",
			"/api/share":     "Returns a list of files contained in the sharing directory.",
		})
	})

	api.Get("/info", middleware.CacheMiddleware(5), getInfo)
	api.Get("/info/net", getNet)
	api.Get("/info/mem", getMem)
	api.Get("/info/disk", getDisk)
	api.Get("/info/ps", getInfoProcess)
	api.Get("/info/usb", getUsb)
	api.Get("/info/cpu", middleware.CacheMiddleware(5), getCpu)
	api.Get("/info/gpio", getGpioList)

	api.Get("/gpio", getGpio)
	api.Get("/gpio/all", middleware.CacheMiddleware(1), getGpioAll)
	api.Patch("/gpio/:pin", updateGpio)

	api.Get("/share", getShare)
	api.Get("/share/*", getShareFile)
	api.Delete("/share/*", middleware.CheckAuth, deleteShareFile)

	api.Delete("/ps/:pid", middleware.CheckAuth, killProcess)
	api.Get("/ps/:pid", getProcessByPid)
}
