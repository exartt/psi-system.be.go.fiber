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

		authorization := c.Get("Authorization")

		// Verifique se o cabeçalho está vazio
		if authorization == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		// Extraia o token da string "Bearer"
		tokenString := strings.Replace(authorization, "Bearer ", "", 1)

		// Parse e valide o token
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

		// Verifique se o token é válido e se expirou
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
					"error": "Expired token",
				})
			}
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
					"error": "Expired token",
				})
			}

			c.Locals("userID", claims["sub"])
			c.Locals("psychologistID", claims["psy_id"])
		}
		return c.Next()
	}
}
