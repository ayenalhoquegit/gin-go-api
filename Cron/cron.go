package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	logger "gin-go-api/log"

	"github.com/go-sql-driver/mysql"
	"github.com/robfig/cron/v3"
)

type OutboxData struct {
	Id        int       `json:"id"`
	Message   string    `json:"message"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdateAt  time.Time `json:"updated_at"`
}

type ApiRequest struct {
	Id        int       `json:"id"`
	Message   string    `json:"message"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

var db *sql.DB
var outboxId int

func getOutbox() {
	var outbox []OutboxData
	rows, err := db.Query("SELECT * FROM outbox where status=0  ORDER BY id LIMIT 0,40")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var u OutboxData
		err := rows.Scan(&u.Id, &u.Message, &u.Status, &u.CreatedAt, &u.UpdateAt)
		if err != nil {
			logger.ErrorLogger.Println(err)
		}
		logger.InfoLogger.Println("Id : ", u.Id)
		outbox = append(outbox, u)
	}

	for _, data := range outbox {
		//insert api data
		_, err = db.Exec("INSERT INTO api_request (message) VALUES (?)", data.Message)
		if err != nil {
			logger.ErrorLogger.Println(err)
		}
		// update outbox
		_, err = db.Exec("UPDATE outbox set status=?, updated_at=? where id=?", 1, currentDateTime(), data.Id)
		if err != nil {
			logger.ErrorLogger.Println(err)
		}
		logger.InfoLogger.Println("Id : " + strconv.Itoa(data.Id) + " | Receiver : number | result : 1")
	}

}

func currentDateTime() string {
	location, err := time.LoadLocation("Asia/Dhaka")
	if err != nil {
		fmt.Println("Error loading location:", err)
	}
	// Get the current time
	currentTime := time.Now().In(location)
	// Define the desired layout  datetime
	layout := "2006-01-02 15:04:05"
	// Format the time using the layout
	formattedTime := currentTime.Format(layout)
	return formattedTime

}

func main() {
	cfg := mysql.Config{
		User:                 "root",
		Passwd:               "",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "golang_api",
		AllowNativePasswords: true,
		ParseTime:            true,
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
	c := cron.New()

	//Schedule a cron job to run every 5 seconds
	c.AddFunc("@every 1s", func() {
		outboxId += 1
		fmt.Println("Cron job executed at:", time.Now())
		getOutbox()
		fmt.Println("Outbox id :", outboxId)
		if outboxId == 3 {
			c.Stop()
		}
	})

	c.Start()
	select {}

}
