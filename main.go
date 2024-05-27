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
	m.HandleFunc("/admin/", adminHandle)
	m.HandleFunc("/ad/update/", updateHandle)
	m.HandleFunc("/ad/delete/", deleteHandle)
	m.HandleFunc("/ad/add/", addHandle)
	m.HandleFunc("/data/", getDataHandle)

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
