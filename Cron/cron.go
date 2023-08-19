package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

type OutboxData struct {
	Id     int    `json:"id"`
	Message   string `json:"message"`
	Status  int `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdateAt time.Time `json:"updated_at"`
}

type ApiRequest struct {
	Id     int    `json:"id"`
	Message   string `json:"message"`
	Status  int `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}



var db *sql.DB

func getOutbox(){
	//var outbox []OutboxData
	rows, err := db.Query("SELECT * FROM outbox where status=0  ORDER BY id LIMIT 0,10")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var u OutboxData
		err := rows.Scan(&u.Id, &u.Message, &u.Status, &u.CreatedAt,&u.UpdateAt)
		if err != nil{
			log.Fatal(err)
		}
		//insert api data
		result,err := db.Exec("INSERT INTO api_request (message) VALUES (?)",u.Message)
		println(result)
		if err != nil {
			log.Fatal(err)
		}
		// update outbox
		result, err = db.Exec("UPDATE outbox set status=?, updated_at=? where id=?",1,currentDateTime(), u.Id)
		println(result)
		if err != nil {
			log.Fatal(err) 
		}
		//fmt.Println("created_at ", u.CreatedAt)
		//fmt.Println("updated_at ", u.UpdateAt)
		//outbox = append(outbox, u)
	}
    //fmt.Println("current datetime", currentDateTime())
	
}

func currentDateTime() string{
	location, err := time.LoadLocation("Asia/Dhaka")
	if err != nil {
		fmt.Println("Error loading location:", err)
	}
	// Get the current time in UTC
	currentTime := time.Now().In(location)

	// Define the desired layout for UTC datetime
	layout := "2006-01-02 15:04:05"

	// Format the UTC time using the layout
	formattedTime := currentTime.Format(layout)

	return formattedTime;
	
}

// func init() {
// 	log.SetLevel(log.InfoLevel)
// 	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
// }

func main() {
	cfg := mysql.Config{
		User:                 "root",
		Passwd:               "",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "golang_api",
		AllowNativePasswords: true,
		ParseTime: true,
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
	getOutbox()
	fmt.Println("Operation success!")
	//c := cron.New()

	// Schedule a cron job to run every 5 seconds
    // c.AddFunc("@every 1s", func() {
		
	// })

	// c.Start()

	// // Keep the program running to execute cron jobs
	// select {}

	

	// 	fmt.Println("Cron job executed at:", time.Now())

	// // Funcs may also be added to a running Cron
	// log.Info("Add new job to a running cron")
	// entryID2, _ := c.AddFunc("*/2 * * * *", func() { log.Info("[Job 2]Every two minutes job\n") })
	// printCronEntries(c.Entries())
	// time.Sleep(5 * time.Minute)

	// //Remove Job2 and add new Job2 that run every 1 minute
	// log.Info("Remove Job2 and add new Job2 with schedule run every minute")
	// c.Remove(entryID2)
	// c.AddFunc("*/1 * * * *", func() { log.Info("[Job 2]Every one minute job\n") })
	// time.Sleep(5 * time.Minute)
	// const sliceSize = 10000
	// for i := 0; i < sliceSize; i++ {
	// 	db.Exec("INSERT INTO outbox (message) VALUES (?)", "message"+strconv.Itoa(i))
		
	// }
	// fmt.Println("Successfully inserted")



	

}

func printCronEntries(cronEntries []cron.Entry) {
	log.Infof("Cron Info: %+v\n", cronEntries)
}