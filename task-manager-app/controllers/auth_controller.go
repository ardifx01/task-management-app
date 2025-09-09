package controllers

import (
	"task-manager-app/config"
	"task-manager-app/models"
	"task-manager-app/utils"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var validate = validator.New()

func Register(c *fiber.Ctx)error{
	var user models.User

	if err := c.BodyParser(&user);err!=nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error" : "invalid parsing body",
		})
	}
	

	if err := validate.Struct(&user);err!=nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error" : err.Error(),
		})
	}

	collection := config.DB.Collection("users")
	var existingUser models.User
	err := collection.FindOne(c.Context(), bson.M{"email" : user.Email}).Decode(&existingUser)
	if err ==nil{
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error" : "Email sudah digunakan",
		})
	}
	hashedPassword,err := utils.HashPassword(user.Password)
	if err !=nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error" : "Gagal hashing password",
		})
	}
	user.ID = primitive.NewObjectID()
	user.Password = hashedPassword
	user.CreatedAt = time.Now()

	result,err := collection.InsertOne(c.Context(), user)
	if err !=nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error" : err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message" : "user berhasil di registrasi", "userId" : result.InsertedID,
	})
}
func Login(c *fiber.Ctx)error{
	var body struct{
		Email string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	if err := c.BodyParser(&body);err!=nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error" : "invalid parse body",
		})
	}

	if err := validate.Struct(&body);err!=nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err" : err.Error(),
		})
	}

	collection := config.DB.Collection("users")
	var user models.User
	err := collection.FindOne(c.Context(), bson.M{"email": body.Email}).Decode(&user)
	
	if err!=nil{
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error" : err.Error(),
		})
	}
	if !utils.CheckPasswordHash(body.Password, user.Password){
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error" : "invalid credential",
		})
	}
	token, err := utils.GenerateToken(user.ID)
	if err!=nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error" : "tidak bisa generate token",
		})
	}
	return c.JSON(fiber.Map{"token" : token})
}
