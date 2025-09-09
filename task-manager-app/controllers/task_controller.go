package controllers

import (
	"task-manager-app/config"
	"task-manager-app/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskRequest struct {
    Title       string    `json:"title" validate:"required"`
    Description string    `json:"description"`
    Priority    string    `json:"priority" validate:"oneof=low medium high"`
    DueDate     time.Time `json:"dueDate,omitempty"`
}

func CreateTask(c *fiber.Ctx) error {
    projectID, err := primitive.ObjectIDFromHex(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid Project ID",
        })
    }

    userID, ok := c.Locals("userId").(primitive.ObjectID)
    if !ok {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "User ID not found",
        })
    }

    
    var project models.Project
    projectCollection := config.DB.Collection("projects")
    err = projectCollection.FindOne(
        c.Context(),
        bson.M{"_id": projectID, "ownerId": userID},
    ).Decode(&project)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Project Not found",
        })
    }

    
    var req TaskRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid Request",
        })
    }

    task := models.Task{
        ID:          primitive.NewObjectID(), 
        ProjectID:   projectID,               
        Title:       req.Title,
        Description: req.Description,
        Status:      "todo", 
        Priority:    req.Priority,
        DueDate:     req.DueDate,
        CreatedAt:   time.Now(),
    }

    if err := validate.Struct(&task); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    taskCollection := config.DB.Collection("tasks")
    result, err := taskCollection.InsertOne(c.Context(), task)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Gagal Buat task",
        })
    }

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "message": "Task berhasil dibuat",
        "taskId":  result.InsertedID,
    })
}

func GetTasks(c *fiber.Ctx)error{
	projectID,err := primitive.ObjectIDFromHex(c.Params("id"))
	if err!=nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error" : "Invalid project ID",
		})
	}
	userID,ok := c.Locals("userId").(primitive.ObjectID)
	if !ok{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error" : "User Id not found",
		})
	}
	var project models.Project

	projectCollection := config.DB.Collection("projects")
	err = projectCollection.FindOne(c.Context(), bson.M{"_id" : projectID, "ownerId" : userID}).Decode(&project)
	if err!=nil{
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error" : "Project Not found",
		})
	}
	taskCollection := config.DB.Collection("tasks")
	cursor, err := taskCollection.Find(c.Context(), bson.M{"projectId" : projectID})
	if err !=nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error" : "tidak dapat fetch tasks",
		})
	}
	defer cursor.Close(c.Context())

	var tasks []models.Task
	if err = cursor.All(c.Context(), &tasks);err!=nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error" : "gagal decode tasks",
		})
	}
	return c.JSON(tasks)
}

func GetTask(c *fiber.Ctx)error{
	taskID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err!=nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error" : "Tidak dapat Menemukan ID",
		})
	}

	taskCollection := config.DB.Collection("tasks")
	var task models.Task
	err = taskCollection.FindOne(c.Context(), bson.M{"_id" : taskID}).Decode(&task)
	if err!=nil{
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error" : "Task Not Found",
		})
	}

	userID, ok := c.Locals("userId").(primitive.ObjectID)
	if !ok{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error" : "User Not Found",
		})
	}
	var project models.Project
	projectCollection := config.DB.Collection("projects")
	err = projectCollection.FindOne(c.Context(), bson.M{"_id" : task.ProjectID, "ownerId" : userID}).Decode(&project)
	if err!=nil{
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error" : "Kamu tidak punya akses untuk mengakses task ini",
		})
	}
	return c.JSON(task)
}

func UpdateTask(c *fiber.Ctx)error{
	taskID,err := primitive.ObjectIDFromHex(c.Params("id"))
	if err!=nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error" : "Invalid Task ID",
		})
	}
	taskCollection := config.DB.Collection("tasks")
	var tasks models.Task
	err = taskCollection.FindOne(c.Context(), bson.M{"_id" : taskID}).Decode(&tasks)
	if err!=nil{
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error" : "Task Not Found",
		})
	}
	userID,ok := c.Locals("userId").(primitive.ObjectID)
	if !ok{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error" : "User ID not found",
		})
	}
	var project models.Project
	projectCollection := config.DB.Collection("projects")
	err = projectCollection.FindOne(c.Context(), bson.M{"_id" : tasks.ProjectID, "ownerId" : userID}).Decode(&project)
	if err!=nil{
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error" : "Kamu tidak punya akses untuk task ini",
		})
	}
	var updatedData map[string]interface{}

	if err := c.BodyParser(&updatedData);err!=nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error" : "Invalid parsing Body",
		})
	}
	if _,ok := updatedData["status"];ok && (updatedData["status"] != "todo" && updatedData["status"] != "in-progress" && updatedData["status"] != "done"){
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error" : "Invalid status Value",
		})
	} 
	if _,ok := updatedData["priority"];ok &&(updatedData["priority"] != "low" && updatedData["priority"] != "medium" && updatedData["priority"] != "high"){
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error" : "Invalid Priority Value",
		})
	}
	update := bson.M{"$set" : updatedData}
	result,err := taskCollection.UpdateOne(c.Context(), bson.M{"_id" : taskID}, update)

	if err !=nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error" : "Gagal update task",
		})
	}
	if result.ModifiedCount == 0{
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error" : "tidak ada task yang diupdate",
		})
	}
	return c.JSON(fiber.Map{
		"message" : "Task berhasil diupdate",
	})
}

func DeleteTask(c *fiber.Ctx)error{
	taskID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err!=nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error" : "Invalid Task ID",
		})
	}
	taskCollection := config.DB.Collection("tasks")
	var task models.Task
	err = taskCollection.FindOne(c.Context(), bson.M{"_id" : taskID}).Decode(&task)
	if err!=nil{
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error" : "Task Not Found",
		})
	}
	userID, ok := c.Locals("userId").(primitive.ObjectID)
	if !ok{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error" : "User ID is not found",
		})
	}
	var project models.Project
	projectCollection := config.DB.Collection("projects")
	err = projectCollection.FindOne(c.Context(), bson.M{"_id" : task.ProjectID, "ownerId" : userID}).Decode(&project)
	if err!=nil{
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error" : "Kamu tidak punya akses ke task ini",
		})
	}
	result, err := taskCollection.DeleteOne(c.Context(), bson.M{"_id" : taskID})
	if err!=nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error" : "Tidak bisa Hapus task",
		})
	}
	if result.DeletedCount ==0{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error" : "Task tidak ditemukan",
		})
	}

	return c.JSON(fiber.Map{
		"message" : "Task berhasil dihapus",
	})
}


func FilterTask(c *fiber.Ctx)error{
	userID,ok := c.Locals("userId").(primitive.ObjectID)
	if !ok{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error" : "User ID not found",
		})
	}
	projectCollection := config.DB.Collection("projects")
	projectCursor, err := projectCollection.Find(c.Context(), bson.M{"_id" : userID})
	if err!=nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error" : "Gagal mengambil projek user",
		})
	}
	defer projectCursor.Close(c.Context())

	var projects []models.Project
	if err = projectCursor.All(c.Context(), &projects);err!=nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error" : "Gagal decode Projek user",
		})
	}

	if len(projects) ==0{
		return c.JSON([]models.Task{})
	}
	projectIDs := make([]primitive.ObjectID, len(projects))
	for i,p := range projects{
		projectIDs[i] = p.ID
	}	
	filter := bson.M{
		"projectId" : bson.M{
			"$in" : projectIDs,
		},
	}
	queryStatus := c.Query("status")
	if queryStatus != ""{
		filter["status"] = queryStatus
	}
	queryPriority := c.Query("priority")
	if queryPriority != ""{
		filter["priority"] = queryPriority
	}
	queryDueBefore := c.Query("dueBefore")
	if queryDueBefore != ""{
		dueDate,err := time.Parse(time.RFC3339, queryDueBefore)
		if err==nil{
				filter["dueDate"] = bson.M{"$lte": dueDate}
		}
	}
	taskCollection := config.DB.Collection("tasks")
	taskCursor, err := taskCollection.Find(c.Context(), filter)
	if err != nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error" : "Gagal Fetch task",
		})
	}
	defer taskCursor.Close(c.Context())
	var tasks []models.Task
	if err = taskCursor.All(c.Context(), &tasks);err!=nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error" : "Gagal decode task",
		})
	}
	return c.JSON(tasks)
}