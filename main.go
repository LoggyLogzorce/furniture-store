package main

import (
	"furniture_store/db"
	"log"
	"net/http"
	"time"
)

func main() {
	m := http.NewServeMux()

	m.HandleFunc("/", userHandle)
	m.HandleFunc("/ad/", updateHandle)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      m,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	db.Connect()
	db.Migrate()

	log.Println("Listening port :8080...")
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
