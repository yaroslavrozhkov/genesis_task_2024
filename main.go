package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func main() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	DB.AutoMigrate(&Subscriber{})

	router := mux.NewRouter()
	router.HandleFunc("/api/subscribe", SubscribeHandler).Methods("POST")
	router.HandleFunc("/api/rate", CurrencyHandler).Methods("GET")

	go StartScheduler()

	log.Println("Server started at :8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
