package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/m0rk0vka/avito-tech-backend-trainee-assigment-2023/router"
)

func main() {
	r := router.Router()
	fmt.Println("Starting server on port 8080...")

	log.Fatal(http.ListenAndServe(":8080", r))
}
