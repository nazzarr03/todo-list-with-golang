package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nazzarr03/todo-list-with-golang/controller"
)
func main() {
	r := gin.Default()

	r.GET("/todos", controller.GetTodos)
	r.GET("/todos/:id", controller.GetTodo)
	r.POST("/todos", controller.CreateTodo)
	r.PUT("/todos/:id", controller.UpdateTodo)
	r.DELETE("/todos/:id", controller.DeleteTodo)
	r.PATCH("/todos/:id/complete", controller.CompleteTodo)

	r.Run(":8080")
}