package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Support struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID     string             `bson:"user_id" json:"user_id"`
	Issue      string             `bson:"issue" json:"issue"`
	Status     string             `bson:"status" json:"status"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	ResolvedAt time.Time          `bson:"resolved_at" json:"resolved_at"`
}
