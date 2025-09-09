package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ProjectID primitive.ObjectID `bson:"projectId,omitempty" json:"projectId,omitempty"`
	Title string	`bson:"title" json:"title" validate:"required"`
	Description string `bson:"description" json:"description"`
	Status string `bson:"status" json:"status" validate:"oneof=todo in-progress done"`
	Priority string `bson:"priority" json:"priority" validate:"oneof=low medium high"`
	DueDate time.Time `bson:"dueDate" json:"dueDate,omitempty"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt,omitempty"`
}