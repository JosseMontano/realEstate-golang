package controllers

import (
	"fmt"
	"time"

	"github.com/JosseMontano/estateInTheCloud/database"
	"github.com/JosseMontano/estateInTheCloud/middleware"
	"github.com/JosseMontano/estateInTheCloud/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

func ValidateStruct(user models.User) []*models.ErrorResponse {
	var errors []*models.ErrorResponse
	err := validate.Struct(user)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element models.ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

func Register(c *fiber.Ctx) error {

	var data map[string]string

	urlPhot :=
		"https://res.cloudinary.com/dny08tnju/image/upload/v1672280689/real_estates/desconocido_hgz7m2.jpg"

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	if data["password"] != data["password_confirm"] {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "password do not match",
		})
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	user := models.User{
		UserName:        data["username"],
		Email:           data["email"],
		CellphoneNumber: data["cellphone_number"],
		Password:        password,
		UrlPhoto:        urlPhot,
	}

	errors := ValidateStruct(user)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": errors,
		})
	}

	result := database.DB.Create(&user)

	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(result.Error)
	}

	return c.JSON(user)
}

func SingIn(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User

	database.DB.Where("email=?", data["email"]).First(&user)

	if user.Id == 0 {
		c.Status(404)
		return c.JSON(fiber.Map{
			"message": "not found",
		})
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "incorrect password",
		})
	}

	timeExp := time.Now().Add(24 * time.Hour)
	tokenString, err := middleware.GenerateJwt(user, timeExp)

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
/* 
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Expires:  timeExp,
		HTTPOnly: true,
	}

	c.Cookie(&cookie) */

	return c.JSON(fiber.Map{
		"message": "success",
		"token":   tokenString,
	})
}

func User(c *fiber.Ctx) error {
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

	return c.JSON(claims)

}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})

}
