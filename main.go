package main

import (
	"log"
	"os"

	"image-api/redis"
	"image-api/utils"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
)

func startServer() {
	e := echo.New()
	e.POST("/upload", utils.UploadImage)
	e.GET("/download", utils.SendImage)
	e.File("/start", "index.html")

	err := e.Start(":8000")
	if err != nil {
		panic(err)
	}
}

func processENV() (storagePath string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	storagePath = os.Getenv("STORAGE_PATH")
	return storagePath
}

func main() {
	utils.StoragePath = processENV()

	log.Println("service started")

	redis.SetupRedisClient()

	utils.StartRabbitMQ()

	go utils.RMQ.RecieveMessages()

	startServer()
}
