package config

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strings"
	"time"
)

func JWTMiddleware() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		authorization := c.Cookies("Authorization")
		if authorization == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		tokenString := strings.Replace(authorization, "Bearer ", "", 1)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(os.Getenv("SECRET")), nil
		})

		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
					"error": "expired token",
				})
			}
			c.Locals("userID", claims["sub"])

			c.Cookie(&fiber.Cookie{
				Name:     "Authorization",
				Value:    tokenString,
				HTTPOnly: true,
			})
		}
		return c.Next()
	}
}
