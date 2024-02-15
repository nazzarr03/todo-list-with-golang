package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TodoItem struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title      string             `json:"title"`
	Desc       string             `json:"desc"`
	Completed  bool               `json:"completed"`
	Created_at time.Time          `json:"created_at"`
	Updated_at time.Time          `json:"updated_at"`
}

type TodoList struct {
	Todos []TodoItem
}
