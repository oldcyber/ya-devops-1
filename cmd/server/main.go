package main

import (
	"log"
	"net/http"

	"ya-devops-1/internal/server"
)

func main() {
	http.HandleFunc("/", server.GetRoot)
	http.HandleFunc("/update/", server.GetMetrics)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
