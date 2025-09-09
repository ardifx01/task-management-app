package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Project struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    OwnerID     primitive.ObjectID `bson:"ownerId,omitempty" json:"ownerId,omitempty"`
    Name        string             `bson:"name" json:"name" validate:"required"`
    Description string             `bson:"description" json:"description"`
    CreatedAt   time.Time          `bson:"createdAt" json:"createdAt,omitempty"`
}
