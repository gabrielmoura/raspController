package routes

import (
	"github.com/gabrielmoura/raspController/pkg/vchiq"
	"github.com/gofiber/fiber/v2"
	"os/exec"
	"time"
)

// killProcess godoc
// @description Kills a process by PID.
// @tags info
// @url /api/info/ps/{pid}
func killProcess(c *fiber.Ctx) error {
	pid := c.Params("pid")
	err := killProcessByPid(pid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Process killed",
	})
}
func killProcessByPid(pid string) error {
	cmd := exec.Command("kill", pid)
	return cmd.Run()
}

// getInfoProcess godoc
// @description Returns all processes and their information.
// @tags info
// @url /api/info/ps
func getInfoProcess(c *fiber.Ctx) error {
	ps, err := vchiq.ListProcesses()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"processes":    ps,
		"count":        len(ps),
		"reading_date": time.Now().Format("2006-01-02 15:04:05"),
	})
}

// getProcessByPid godoc
// @description Returns a process by PID.
// @tags info
// @url /api/info/ps/{pid}
func getProcessByPid(c *fiber.Ctx) error {
	pid := c.Params("pid")
	ps, err := vchiq.GetProcessByPid(pid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"processes":    ps,
		"count":        len(ps),
		"reading_date": time.Now().Format("2006-01-02 15:04:05"),
	})
}
