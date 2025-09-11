package middleware

import (
	"os"
	"strings"

	"github.com/danielreinhard1129/fiber-clean-arch/pkg/exception"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return exception.UnauthorizedError{Message: "Missing or malformed JWT"}
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return exception.UnauthorizedError{Message: "Authorization header format must be Bearer {token}"}
		}

		secret := os.Getenv("JWT_SECRET")
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, exception.UnauthorizedError{Message: "Unexpected signing method"}
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return exception.UnauthorizedError{Message: "Invalid or expired JWT"}
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Locals("id", claims["id"])
			c.Locals("email", claims["email"])
			c.Locals("role", claims["role"])
		}

		return c.Next()
	}
}
