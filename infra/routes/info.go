package routes

import (
	"log"
	"time"

	"github.com/gabrielmoura/raspController/pkg/vchiq"
	"github.com/gofiber/fiber/v2"
)

// getInfo returns system information.
func getInfo(c *fiber.Ctx) error {
	info := make(fiber.Map)
	info["reading_date"] = time.Now().Format("2006-01-02 15:04:05")
	if hostname, err := vchiq.GetHostname(); err != nil {
		log.Println("Error getting hostname:", err)
	} else {
		info["hostname"] = hostname
	}

	if vchiq.IsVcgencmdInstalled() {
		if volt, err := vchiq.GetCoreVolt(); err != nil {
			log.Println("Error getting core voltage:", err)
		} else {
			info["core_voltage"] = volt
		}

		if temp, err := vchiq.GetGPUTemp(); err != nil {
			log.Println("Error getting GPU temperature:", err)
		} else {
			info["gpu_temp"] = temp
		}

		if throttled, err := vchiq.GetThrottled(); err != nil {
			log.Println("Error getting throttled status:", err)
		} else {
			info["throttled"] = throttled
		}

		if arm, gpu, err := vchiq.GetMem(); err != nil {
			log.Println("Error getting memory info:", err)
		} else {
			info["arm_mem"] = arm
			info["gpu_mem"] = gpu
		}
	}

	if boot, root, home, err := vchiq.GetDiskUsage(); err != nil {
		log.Println("Error getting disk usage:", err)
	} else {
		info["disk"] = fiber.Map{
			"boot": boot,
			"root": root,
			"home": home,
		}
	}

	if usedPercent, err := vchiq.GetMemoryUsagePercent(); err != nil {
		log.Println("Error getting memory usage percent:", err)
	} else {
		info["memory_percent"] = fiber.Map{
			"used": usedPercent,
		}
	}

	if total, free, used, err := vchiq.GetMemory(); err != nil {
		log.Println("Error getting memory info:", err)
	} else {
		info["memory"] = fiber.Map{
			"total": total,
			"free":  free,
			"used":  used,
		}
	}

	if temp, err := vchiq.GetCPUTemp(); err != nil {
		log.Println("Error getting CPU temperature:", err)
	} else {
		info["cpu_temp"] = temp
	}

	return c.Status(fiber.StatusOK).JSON(info)
}

// getInfoProcess returns all processes and their information.
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
