package middleware

import (
	"github.com/gabrielmoura/raspController/configs"
	"log"
	"time"

	"github.com/gabrielmoura/raspController/infra/db"
	"github.com/gofiber/fiber/v2"
)

// CacheMiddleware godoc
// @description Middleware for caching responses
func CacheMiddleware(ttl int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Method() != fiber.MethodGet {
			return c.Next()
		}

		url := c.OriginalURL()

		// Attempts to retrieve response from storage
		cachedBody, err := db.DB.Get([]byte(url))
		if err == nil {
			return c.Send(cachedBody)
		}

		// If you don't find the answer, follow the original flow
		if err := c.Next(); err != nil {
			return err
		}

		// Get the generated response
		body := c.Response().Body()

		// Stores the response in storage
		err = db.DB.PutWithTTL([]byte(url), body, time.Second*time.Duration(ttl))
		if err != nil {
			log.Println("Error storing cache for", url, err)
		}

		return nil
	}
}

// CheckAuth godoc
// @description Middleware for authentication
func CheckAuth(c *fiber.Ctx) error {
	if c.Get("Authorization") != "Bearer "+configs.Conf.AuthToken {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	return c.Next()
}
