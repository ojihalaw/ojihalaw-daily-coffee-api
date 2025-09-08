package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ojihalawa/daily-coffee-api.git/internal/utils"
)

func AuthMiddleware(jwtMaker *utils.JWTMaker) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).
				JSON(utils.ErrorResponse(fiber.StatusUnauthorized, "missing token"))
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := jwtMaker.VerifyAccessToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).
				JSON(utils.ErrorResponse(fiber.StatusUnauthorized, "invalid or expired token"))
		}

		// inject user info ke context
		c.Locals("userID", claims.UserID)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}
