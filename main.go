package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

type User struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Gender string `json:"gender"`
}

var db *sql.DB

func getAllUser(c *gin.Context) {
	var user []User
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var u User
		rows.Scan(&u.Id, &u.Name, &u.Email, &u.Gender)
		user = append(user, u)
	}
	c.IndentedJSON(http.StatusOK, user)
}
func postUser(c *gin.Context) {
	var newUser User
	if err := c.BindJSON(&newUser); err != nil {
		return
	}
	result, err := db.Exec("INSERT INTO users (name, email, gender) VALUES (?, ?, ?)", newUser.Name, newUser.Email, newUser.Gender)
	if err != nil {
		log.Fatal(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	newUser.Id = int(id)
	c.IndentedJSON(http.StatusCreated, newUser)
}
func findUserById(c *gin.Context) {
	var user User
	userId, _ := strconv.Atoi(c.Param("id"))
	row := db.QueryRow("SELECT * FROM users WHERE id = ?", userId)
	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Gender)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
	}
	c.IndentedJSON(http.StatusOK, user)

}

func removeUserById(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Param("id"))
	result, err := db.Exec("DELETE FROM users WHERE id = ?", userId)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}
	_, err = result.RowsAffected()
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "user deleted"})

}

func main() {

	cfg := mysql.Config{
		User:                 "root",
		Passwd:               "",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "golang_api",
		AllowNativePasswords: true,
	}

	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
	router := gin.Default()
	router.GET("/users", getAllUser)
	router.POST("/users", postUser)
	router.GET("/users/:id", findUserById)
	router.DELETE("/users/:id", removeUserById)
	router.Run("localhost:8080")
}
