package routes

import (
	"github.com/gabrielmoura/raspController/configs"
	"github.com/gabrielmoura/raspController/pkg/files"
	"github.com/gofiber/fiber/v2"
)

// getShare retorna uma lista de arquivos contidos no diretório de compartilhamento.
func getShare(c *fiber.Ctx) error {
	if len(configs.Conf.ShareDir) > 0 {
		files, err := files.ListDirectory(configs.Conf.ShareDir)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"files": files,
		})
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error": "Not configured, set SHARE_DIR",
	})
}

// getShareFile retorna um arquivo do diretório de compartilhamento ou uma lista de arquivos se o caminho for um diretório.
func getShareFile(c *fiber.Ctx) error {
	if len(configs.Conf.ShareDir) > 0 {
		filePath := configs.Conf.ShareDir + "/" + c.Params("*")

		if files.IsFolder(filePath) {
			files, err := files.ListDirectory(filePath)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": err.Error(),
				})
			}

			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"files": files,
			})
		}

		return c.SendFile(filePath)
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error": "Not configured, set SHARE_DIR",
	})
}
