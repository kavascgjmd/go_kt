package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
     ID			primitive.ObjectID `bson:"_id"`
	 Created_at *time.Time		   `jaon:"created_at"`
	 Updated_at *time.Time         `json:"update_at"`
	 Order_id    string            `json:"order_id"`
     Table_id    *string           `json:"table_id"`
}