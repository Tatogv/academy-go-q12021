package main

import (
	"api/router"
	"api/usecase"
	"log"
)

func main() {
	config, err := usecase.LoadConfig()
	if err != nil {
		log.Println("Error loading config. Make sure you have an api.env file with the correct configurations.")
		return
	}
	router.CreateRouter(config)
}
