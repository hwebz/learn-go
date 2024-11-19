package main

import (
	"fmt"
	"strings"

	"github.com/go-jwt-auth/initializers"
	"github.com/go-jwt-auth/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func DeserializeUser(c *fiber.Ctx) error {
	fmt.Println("tokenString")
	var tokenString string
	authorization := c.Get("Authorization")

	if strings.HasPrefix(authorization, "Bearer ") {
		tokenString = strings.TrimPrefix(authorization, "Bearer ")
	} else if c.Cookies("token") != "" {
		tokenString = c.Cookies("token")
	}

	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "failed",
			"message": "Unauthorized",
		})
	}

	config, _ := initializers.LoadEnv(".")

	tokenByte, err := jwt.Parse(tokenString, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %s", jwtToken.Header["alg"])
		}

		return []byte(config.JwtSecret), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "failed",
			"message": fmt.Sprintf("Invalid token: %v", err),
		})
	}

	claims, ok := tokenByte.Claims.(jwt.MapClaims)
	if !ok || !tokenByte.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "failed",
			"message": "Invalid token claims",
		})
	}

	var user models.User
	initializers.DB.First(&user, "id = ?", claims["sub"])

	if user.ID.String() != claims["sub"] {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "failed",
			"message": "The user belonging to this token no longer exists",
		})
	}

	c.Locals("user", models.FilterUserRecord(&user))

	return c.Next()
}

func allowedRoles(allowedRoles []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := c.Locals("user").(models.UserResponse)

		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "failed",
				"message": "Access denied. User not authenticated",
			})
		}

		roleAllowed := false
		for _, allowedRole := range allowedRoles {
			if user.Role == allowedRole {
				roleAllowed = true
				break
			}
		}

		if !roleAllowed {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "failed",
				"message": "Access denied. You are not allowed to perform this action.",
			})
		}

		return c.Next()
	}
}
