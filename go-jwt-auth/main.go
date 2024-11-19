package main

import (
	"fmt"
	"log"

	"github.com/go-jwt-auth/initializers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	env, err := initializers.LoadEnv(".")
	if err != nil {
		log.Fatal("🛑 Could not load environment variables", err)
	}
	initializers.ConnectDB(&env)

	app := fiber.New()
	micro := fiber.New()

	app.Mount("/api", micro)
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     env.ClientOrigin,
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowMethods:     "GET, POST",
		AllowCredentials: true,
	}))

	SetupRoutes(micro)

	micro.Get("/healthchecker", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "JSON Web Token Authentication and Authorization in Golang",
		})
	})

	app.All("*", func(c *fiber.Ctx) error {
		path := c.Path()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": fmt.Sprintf("Path: %v doesn not exists on this server", path),
		})
	})

	log.Fatal(app.Listen(":8000"))
}
