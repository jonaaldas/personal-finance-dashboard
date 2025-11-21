package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jonaaldas/personal-finance-dashboard/plaid"
)

func main() {
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello World",
		})
	})

	app.Post("/api/set_access_token", func(c *fiber.Ctx) error {
		return plaid.GetAccessToken(c)
	})

	app.Post("/api/create_link_token", func(c *fiber.Ctx) error {
		linkToken, err := plaid.LinkTokenCreate(nil)

		if err != nil {
			fmt.Printf("Error creating link token1: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"link_token": linkToken,
		})
	})

	app.Get("/api/accounts", func(c *fiber.Ctx) error {
		return plaid.Accounts(c)
	})

	app.Get("/api/liabilities", func(c *fiber.Ctx) error {
		return plaid.Liabilities(c)
	})

	app.Get("/api/transactions", func(c *fiber.Ctx) error {
		return plaid.Transactions(c)
	})

	app.Get("/api/access_tokens", func(c *fiber.Ctx) error {
		res, err := plaid.GetAllAccessTokens()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(fiber.Map{
			"tokens": res.AccessTokens,
		})
	})

	log.Fatal(app.Listen(":" + getPort()))
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	return port
}
