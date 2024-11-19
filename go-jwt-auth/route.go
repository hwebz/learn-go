package main

import (
	"github.com/go-jwt-auth/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app fiber.Router) {
	authRoutes := app.Group("/auth")
	authRoutes.Post("/register", handlers.SignUpUser)
	authRoutes.Post("/login", handlers.SignInUser)
	authRoutes.Get("/logout", DeserializeUser, handlers.LogOutUser)

	app.Get("/users/me", DeserializeUser, handlers.GetMeHandler)
	app.Get("/users/", DeserializeUser, allowedRoles([]string{"admin", "moderator"}), handlers.GetUsersHandler)
}
