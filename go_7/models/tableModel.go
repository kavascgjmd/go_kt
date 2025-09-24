package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Table struct {
	ID        primitive.ObjectID   `bson:"_id"`
	Table_id  string               `json:"table_id"`
	Created_at *time.Time          `json:"created_at"`
	Updated_at *time.Time          `json:"updated_at"`
}