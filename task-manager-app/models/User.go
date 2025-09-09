package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
    ID        primitive.ObjectID `json:"id" bson:"_id"`
    Username  string             `json:"username" bson:"username" validate:"required,min=3"`
    Email     string             `json:"email" bson:"email" validate:"required,email"`
    Password  string             `json:"password" bson:"password" validate:"required,min=6"`
    CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
}
