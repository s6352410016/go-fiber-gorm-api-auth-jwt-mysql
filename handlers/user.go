package handlers

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/s6352410016/go-fiber-gorm-api-auth-jwt-mysql/database"
	"github.com/s6352410016/go-fiber-gorm-api-auth-jwt-mysql/models"
	"github.com/s6352410016/go-fiber-gorm-api-auth-jwt-mysql/request"
	"golang.org/x/crypto/bcrypt"
)

func CreateToken(u *models.User) (string, string, error) {
	atClaims := jwt.MapClaims{
		"id":       u.ID,
		"fullname": u.FullName,
		"username": u.UserName,
		"email":    u.Email,
		"exp":      time.Now().Add(time.Minute * 5).Unix(),
	}
	rtClaims := jwt.MapClaims{
		"id":       u.ID,
		"fullname": u.FullName,
		"username": u.UserName,
		"email":    u.Email,
		"exp":      time.Now().Add(time.Hour).Unix(),
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)

	accessToken, err := at.SignedString([]byte(os.Getenv("AT_SECRET")))
	if err != nil {
		return "", "", err
	}
	refreshToken, err := rt.SignedString([]byte(os.Getenv("RT_SECRET")))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func SignUp(c *fiber.Ctx) error {
	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request Data",
		})
	}

	if user.FullName == "" || user.UserName == "" || user.Password == "" || user.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Input Is Required",
		})
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error,
		})
	}

	user.Password = string(hashPassword)
	result := database.DB.Create(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": result.Error,
		})
	}

	accessToken, refreshToken, err := CreateToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "55555",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func SignIn(c *fiber.Ctx) error {
	userRequest := new(request.UserRequest)
	user := new(models.User)

	if err := c.BodyParser(userRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request Data",
		})
	}

	if userRequest.UserNameOrEmail == "" || userRequest.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Input Is Required",
		})
	}

	database.DB.Where("user_name = ? OR email = ?", userRequest.UserNameOrEmail, userRequest.UserNameOrEmail).First(&user)
	if user.ID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid Credential",
		})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userRequest.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid Credential",
		})
	}

	accessToken, refreshToken, err := CreateToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error,
		})
	}

	return c.JSON(fiber.Map{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func Profile(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := claims["id"].(float64)
	fullName := claims["fullname"].(string)
	userName := claims["username"].(string)
	email := claims["email"].(string)

	return c.JSON(fiber.Map{
		"id":       id,
		"fullname": fullName,
		"username": userName,
		"email":    email,
	})
}

func Refresh(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := claims["id"].(float64)
	fullName := claims["fullname"].(string)
	userName := claims["username"].(string)
	email := claims["email"].(string)

	userData := new(models.User)
	userData.ID = uint(id)
	userData.FullName = fullName
	userData.UserName = userName
	userData.Email = email

	accessToken, refreshToken, err := CreateToken(userData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error,
		})
	}

	return c.JSON(fiber.Map{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}
