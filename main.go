package main

import (
	"log"
	"net/http"

	"image-api/redis"
	"image-api/utils"
)

func setupEndpoints() {
	http.HandleFunc("/upload", utils.UploadImage)
	http.ListenAndServe(":8000", nil)
}

func main() {
	log.Println("service started")

	redis.SetupRedisClient()

	utils.StartRabbitMQ()

	go utils.RMQ.RecieveMessages()

	setupEndpoints()
}
