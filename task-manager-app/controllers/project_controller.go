package controllers

import (
	
	"task-manager-app/config"
	"task-manager-app/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProjectRequest struct {
    Name        string `json:"name" validate:"required"`
    Description string `json:"description"`
}
func CreateProject(c *fiber.Ctx) error {
    var req ProjectRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid Request Body",
        })
    }

    userID, ok := c.Locals("userId").(primitive.ObjectID)
    if !ok {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "USER ID tidak ditemukan di context",
        })
    }

    project := models.Project{
        ID:          primitive.NewObjectID(), 
        OwnerID:     userID,                  
        Name:        req.Name,
        Description: req.Description,
        CreatedAt:   time.Now(),
    }

    if err := validate.Struct(&project); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    collection := config.DB.Collection("projects")
    result, err := collection.InsertOne(c.Context(), project)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Gagal membuat project",
        })
    }

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "message":   "Projek berhasil dibuat",
        "projectId": result.InsertedID,
    })
}


func GetProjects( c *fiber.Ctx)error{
	userID, ok := c.Locals("userId").(primitive.ObjectID)
	if !ok{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error" : "User ID tidak ditemukan di context",
		})
	}

	collection := config.DB.Collection("projects")
	cursor, err := collection.Find(c.Context(), bson.M{"ownerId" : userID})
	if err!=nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error" : "Tidak dapat Fetch data Project",
		})
	}
	defer cursor.Close(c.Context())

	var project []models.Project
	if err = cursor.All(c.Context(), &project);err!=nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error" : "tidak dapat mendecode project",
		})
	}
	return c.JSON(project)
}

func GetProject(c *fiber.Ctx)error{
	projectID,err := primitive.ObjectIDFromHex(c.Params("id"))
	if err!=nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error" : "Invalid Project ID",
		})
	}
	userID, ok := c.Locals("userId").(primitive.ObjectID)
	if !ok{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error" : "ID tidak ditemukan",
		})
	}

	collection := config.DB.Collection("projects")
	var project models.Project
	err = collection.FindOne(c.Context(), bson.M{"_id" :projectID, "ownerId" : userID}).Decode(&project)
	if err!=nil{
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error" : "Data projek tidak ditemukan",
		})
	}
	return c.JSON(project)
}

func UpdateProject(c *fiber.Ctx)error{
	projectID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err!=nil{
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error" : "Invalid Project ID",
		})
	}
	userID, ok := c.Locals("userId").(primitive.ObjectID)
	if!ok{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error" : "User ID tidak ditemukan",
		})
	}

	var updatedData map[string]interface{}
	if err := c.BodyParser(&updatedData);err!=nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error" : "Invalid Request Body",
		})
	}

	delete(updatedData, "ownerId")

	update := bson.M{"$set" : updatedData}
	collection := config.DB.Collection("projects")
	result,err := collection.UpdateOne(c.Context(), bson.M{"_id" : projectID, "ownerId" : userID}, update)
	if err!=nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error" : "Gagal Update project",
		})
	}
	if result.ModifiedCount == 0{
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error" : "Project not found",
		})
	}
	return c.JSON(fiber.Map{
		"Message" : "Project berhasil diupdate",
	})
}

func DeleteProject(c *fiber.Ctx)error{
	projectID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err!=nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error" : "Invalid Project ID",
		})
	}
	userID, ok := c.Locals("userId").(primitive.ObjectID)
	if !ok{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error" : "User ID is not found",
		})
	}

	taskCollection := config.DB.Collection("tasks")
	_,err = taskCollection.DeleteMany(c.Context(), bson.M{"projectId":projectID})
	if err !=nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error" : "Gagal menghapus tasks yang berelasi",
		})
	}

	projectCollection := config.DB.Collection("projects")
	result,err := projectCollection.DeleteOne(c.Context(), bson.M{"_id" : projectID, "ownerId" : userID})
	if err!=nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error" : "Gagal delete Projek",
		})
	}
	if result.DeletedCount == 0{
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error" : "Project not found",
		})
	}
	return c.JSON(fiber.Map{
		"Message" : "Projek dan task yang berelasi berhasil dihapus",
	})
}