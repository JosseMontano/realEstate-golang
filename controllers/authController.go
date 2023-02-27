package controllers

import (
	"fmt"
	"time"
	"github.com/JosseMontano/estateInTheCloud/database"
	"github.com/JosseMontano/estateInTheCloud/middleware"
	"github.com/JosseMontano/estateInTheCloud/models"
	"github.com/JosseMontano/estateInTheCloud/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
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

	/* 	if data["password"] != data["password_confirm"] {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "password do not match",
		})
	} */

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

	//Token

	timeExp := time.Now().Add(24 * time.Hour)
	tokenString, err := middleware.GenerateJwt(user, timeExp)

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"auth":  true,
		"token": tokenString,
	})
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
		"auth":  true,
		"token": tokenString,
	})
}

func User(c *fiber.Ctx) error {
	/* 	cookie := c.Cookies("jwt") */
	cookie := c.Get("authorization")

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

	//for of a map(object)
	var userToken interface{}

	for key, element := range claims {
		if key == "user" {
			userToken = element
		}
	}

	return c.JSON(userToken)

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

func SendCodeToGmail(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	numberRandom := utils.Random(5)

	user := models.User{
		Email: data["email"],
	}

	database.DB.Model(&user).
		Where("email = ?", data["email"]).
		Update("CodeRecuperation", numberRandom)

	msg := gomail.NewMessage()
	text := "<b>This is the code: " + numberRandom + "</b>"
	msg.SetHeader("From", "eljosema505@gmail.com")
	msg.SetHeader("To", data["email"])
	msg.SetHeader("Subject", "Recuperate account")
	msg.SetBody("text/html", text)
	/* 	msg.Attach("/home/User/cat.jpg") */

	gmail := utils.DotEnvVariable("GMAIL")
	codeGmail := utils.DotEnvVariable("CODE_GMAIL")
	n := gomail.NewDialer("smtp.gmail.com", 587, gmail, codeGmail)

	if err := n.DialAndSend(msg); err != nil {
		panic(err)
	}

	c.Status(200)
	return c.JSON(fiber.Map{
		"operation": true,
	})
}

func ChangePassword(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	user := models.User{}

	database.DB.
		Where("email = ?", data["email"]).Where("code_recuperation = ?", data["code_recuperation"]).
		First(&user)

	if user.Id == 0 {
		c.Status(404)
		return c.JSON(fiber.Map{
			"message": "the code of the email is incorrect",
		})
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	newUser := models.User{
		Password: password,
	}

	database.DB.Model(&newUser).
		Where("id = ?", user.Id).
		Update("password", password)

	c.Status(200)
	return c.JSON(fiber.Map{
		"operation": true,
	})

}
