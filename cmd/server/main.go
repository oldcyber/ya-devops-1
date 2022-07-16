package main

import (
	"log"
	"net/http"

	"ya-devops-1/internal/server"
)

func main() {
	server.StoredData = make(map[string]server.StoredType)
	http.HandleFunc("/", server.GetRoot)           // метод GET
	http.HandleFunc("/update/", server.GetMetrics) // метод POST
	http.HandleFunc("/value/", server.GetValue)    // метод GET

	log.Fatal(http.ListenAndServe(":8080", nil))
}
