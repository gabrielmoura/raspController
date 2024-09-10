package routes

import (
	"github.com/gabrielmoura/raspController/configs"
	"github.com/gabrielmoura/raspController/pkg/files"
	"github.com/gofiber/fiber/v2"
)

// getShare godoc
// @description Returns a list of files contained in the share directory.
// @tags share
// @url /api/share
func getShare(c *fiber.Ctx) error {
	if len(configs.Conf.ShareDir) > 0 {
		listFiles, err := files.ListDirectory(configs.Conf.ShareDir)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"files": listFiles,
		})
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error": "Not configured, set SHARE_DIR",
	})
}

// getShareFile godoc
// @Description Returns a file from the share directory.
// @tags share
// @url /api/share/*
func getShareFile(c *fiber.Ctx) error {
	if len(configs.Conf.ShareDir) > 0 {
		filePath := configs.Conf.ShareDir + "/" + c.Params("*")

		if files.IsFolder(filePath) {
			listFiles, err := files.ListDirectory(filePath)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": err.Error(),
				})
			}

			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"files": listFiles,
			})
		}

		return c.SendFile(filePath)
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error": "Not configured, set SHARE_DIR",
	})
}

// deleteShareFile godoc
// @description Deletes a file from the share directory.
// @tags share
// @url /api/share/*
func deleteShareFile(c *fiber.Ctx) error {
	if len(configs.Conf.ShareDir) > 0 {
		filePath := configs.Conf.ShareDir + "/" + c.Params("*")

		err := files.DeleteFile(filePath)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "File deleted",
		})
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error": "Not configured, set SHARE_DIR",
	})
}
