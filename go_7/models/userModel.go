package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID			 primitive.ObjectID			`bson:"_id"` 
	User_id      string                     `json:"user_id" validate:"required"`
	FirstName    string                     `json:"first_name" validate:"required"`
	LastName     string                     `json:"last_name" validate:"required"`
	Email        string                     `json:"email" validate:"required"`
    HashPassword string                     `json:"hash_password" validate:"required"`
	Created_at   *time.Time                 `json:"created_at"`
	Updated_at   *time.Time                 `json:"updated_at"`
	RefreshToken *string                    `json:"refresh_token"`
}