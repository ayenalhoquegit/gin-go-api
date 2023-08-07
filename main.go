package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	Id     int
	Name   string
	Email  string
	Gender string
}

var user = []User{
	{Id: 1, Name: "A", Email: "a@a.com", Gender: "male"},
	{Id: 2, Name: "B", Email: "b@b.com", Gender: "female"},
	{Id: 3, Name: "C", Email: "c@c.com", Gender: "male"},
}

func getAllUser(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, user)
}

func main() {
	router := gin.Default()
	router.GET("/users", getAllUser)
	router.Run("localhost:8080")
	fmt.Println("All users:", user)
}
