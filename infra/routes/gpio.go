package routes

import (
	"github.com/gabrielmoura/raspController/infra/gpio"
	"github.com/gabrielmoura/raspController/internal/dto"
	"github.com/gofiber/fiber/v2"
)

// getGpio returns all GPIOs status
func getGpio(c *fiber.Ctx) error {
	if !gpio.CheckChip() {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "GPIO chip not initialized",
		})
	}
	list, err := gpio.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(list)
}

// updateGpio updates a GPIO status
func updateGpio(c *fiber.Ctx) error {
	// Verifique se o chip GPIO está inicializado.
	if !gpio.CheckChip() {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "GPIO chip not initialized",
		})
	}

	// Parse do corpo da requisição para obter o modo do pino.
	var pinMode dto.PinMode
	err := c.BodyParser(&pinMode)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	pinMode.Pin, err = c.ParamsInt("pin")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if err := pinMode.Validation(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Defina o valor do pino no chip GPIO.
	if _, err = gpio.SetBool(pinMode.Pin, pinMode.Direction, pinMode.Value); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(pinMode)
}
