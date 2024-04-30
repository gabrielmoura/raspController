package routes

import (
	"github.com/gabrielmoura/raspController/infra/gpio"
	"github.com/gabrielmoura/raspController/internal/dto"
	"github.com/gofiber/fiber/v2"
)

// getGpio godoc
// @description Returns the status of all configured GPIO pins.
// @tags gpio
// @url /api/gpio
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

// updateGpio godoc
// @description Updates the status of a GPIO pin.
// @tags gpio
// @url /api/gpio/{pin}
func updateGpio(c *fiber.Ctx) error {
	// Check whether the GPIO chip is initialized.
	if !gpio.CheckChip() {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "GPIO chip not initialized",
		})
	}

	// Parse the request body to get the pin mode.
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

	// Set the pin value on the GPIO chip.
	if _, err = gpio.SetBool(pinMode.Pin, pinMode.Direction, pinMode.Value); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(pinMode)
}
