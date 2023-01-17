package middleware

import (
	"fmt"
	"github.com/JosseMontano/estateInTheCloud/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

const SecretKey = "secret"

func GenerateJwt(user models.User, timeExp time.Time) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user"] = user
	claims["exp"] = timeExp.Unix()

	return token.SignedString([]byte(SecretKey))
}

func ValidateJwt(c *fiber.Ctx) error {
	/* 	cookie := c.Cookies("jwt") */
	cookie := c.Get("Token")
	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unauthenticated")
		}
		return []byte("secret"), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Invalid token",
		})
	}

	claims := token.Claims.(jwt.MapClaims)
	exp := claims["exp"].(float64)

	if int64(exp) < time.Now().Local().Unix() {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "token expired",
		})
	}

	return c.Next()
	/* 	return claims, err */

}
