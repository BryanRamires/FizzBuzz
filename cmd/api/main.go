package main

import (
	"log"
	"net/http"

	"github.com/BryanRamires/FizzBuzz/internal/httpapi"
)

func main() {
	router := httpapi.NewRouter()

	log.Println("listening on :8090")
	if err := http.ListenAndServe(":8090", router); err != nil {
		log.Fatal(err)
	}
}
