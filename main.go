package main

import (
	"log"
	"os"

	"image-api/queue"
	"image-api/resizer"
	"image-api/server"
	"image-api/storage"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
)

func startServer() {
	e := echo.New()
	e.POST("/upload", server.HandleUploadImage)
	e.GET("/download", server.HandleSendImage)
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
	resizer.StoragePath = processENV()

	log.Println("service started")

	storage := storage.NewStorageClient(storage.NewRedisClient())

	resizer := resizer.Resize{}

	queue := queue.NewQueueClient(queue.NewRabbitMQ())

	server.InitializeModules(storage, resizer, queue)

	go server.HandleMessages()

	startServer()
}
