package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type User struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Gender string `json:"gender"`
}

var user = []User{
	{Id: 1, Name: "A", Email: "a@a.com", Gender: "male"},
	{Id: 2, Name: "B", Email: "b@b.com", Gender: "female"},
	{Id: 3, Name: "C", Email: "c@c.com", Gender: "male"},
}

func getAllUser(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, user)
}
func postUser(c *gin.Context) {
	var newUser User
	if err := c.BindJSON(&newUser); err != nil {
		return
	}
	user = append(user, newUser)
	c.IndentedJSON(http.StatusCreated, newUser)
}
func findUserById(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Param("id"))
	for _, u := range user {
		if u.Id == userId {
			c.IndentedJSON(http.StatusOK, u)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})

}

func removeUserByIndex(c *gin.Context) {
	index, _ := strconv.Atoi(c.Param("index"))
	fmt.Println("user length ", len(user))
	if len(user) >= index {
		user = append(user[:index], user[index+1:]...)
		c.IndentedJSON(http.StatusOK, gin.H{"message": "user deleted"})
		return
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})

}

func main() {
	router := gin.Default()
	router.GET("/users", getAllUser)
	router.POST("/users", postUser)
	router.GET("/users/:id", findUserById)
	router.DELETE("/users/:index", removeUserByIndex)
	router.Run("localhost:8080")
	fmt.Println("All users:", user)
}
