package middleware

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"net/http"
	"psi-system.be.go.fiber/config"
	"strings"
	"time"
)

func JWTMiddleware() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		authorization := c.Cookies("Authorization")
		if authorization == "" {
			logrus.Warn("Token is necessary to proceed")
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		tokenString := strings.Replace(authorization, "Bearer ", "", 1)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return config.Secret(), nil
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
		}

		return c.Next()
	}
}
