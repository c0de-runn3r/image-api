package main

import (
	"log"

	"image-api/queue"
	"image-api/resizer"
	"image-api/server"
	"image-api/storage"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/labstack/echo"
)

var CFG Config

type Config struct {
	StoragePath     string `env:"STORAGE_PATH" env-default:"image"`
	RedisAddress    string `env:"REDIS" env-default:"localhost:6379"`
	RabbitMQAddress string `env:"RABBIT_MQ" env-default:"amqp://guest:guest@localhost:5672/"`
}

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

func main() {
	err := cleanenv.ReadEnv(&CFG)
	if err != nil {
		log.Println("error reading ENV: ", err)
	}

	resizer.StoragePath = CFG.StoragePath

	log.Println("service started")

	storage := storage.NewStorageClient(storage.NewRedisClient(CFG.RedisAddress))

	resizer := resizer.Resize{}

	queue := queue.NewQueueClient(queue.NewRabbitMQ(CFG.RabbitMQAddress))

	server.InitializeModules(storage, resizer, queue)

	go server.HandleMessages()

	startServer()
}
