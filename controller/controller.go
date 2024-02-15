package controller

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nazzarr03/todo-list-with-golang/database"
	"github.com/nazzarr03/todo-list-with-golang/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetTodos is a function that returns all todos
// I use the mongodb for the database
func GetTodos(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	collection := database.OpenCollection(database.Client, "todo")

	cursor, err := collection.Find(ctx, bson.M{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while getting todos"})
		return
	}

	defer cursor.Close(ctx)

	var todos []bson.M

	if err = cursor.All(ctx, &todos); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while getting todos"})
		return
	}

	c.JSON(http.StatusOK, todos)
}

func GetTodo(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	collection := database.OpenCollection(database.Client, "todo")

	id, _ := primitive.ObjectIDFromHex(c.Param("id"))

	var todo models.TodoItem
	
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&todo)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	c.JSON(http.StatusOK, todo)
}

func CreateTodo(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	collection := database.OpenCollection(database.Client, "todo")

	var todo models.TodoItem

	if err := c.BindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	todo.Created_at = time.Now()
	todo.Updated_at = time.Now()

	result, err := collection.InsertOne(ctx, todo)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while creating todo"})
		return
	}

	todo.ID = result.InsertedID.(primitive.ObjectID)

	c.JSON(http.StatusCreated, todo)
}

func DeleteTodo(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	collection := database.OpenCollection(database.Client, "todo")

	id, _ := primitive.ObjectIDFromHex(c.Param("id"))

	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	
	if err != nil || result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
}

func UpdateTodo(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	collection := database.OpenCollection(database.Client, "todo")

	var todo models.TodoItem

	if err := c.BindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	id, _ := primitive.ObjectIDFromHex(c.Param("id"))

	todo.Updated_at = time.Now()

	result, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": todo})

	if err != nil || result.ModifiedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	c.JSON(http.StatusOK, todo)
}

func CompleteTodo(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	collection := database.OpenCollection(database.Client, "todo")

	id, _ := primitive.ObjectIDFromHex(c.Param("id"))

	result, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"completed": true, "updated_at": time.Now()}})

	if err != nil || result.ModifiedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	var updatedTodo models.TodoItem

	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&updatedTodo)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while getting todo"})
		return
	}

	c.JSON(http.StatusOK, updatedTodo)
}