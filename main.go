package main

import (
	"log"

	"image-api/redis"
	"image-api/utils"

	"github.com/labstack/echo"
)

func setupEndpoints() {
	e := echo.New()
	e.POST("/upload", utils.UploadImage)
	e.GET("/download", utils.SendImage)
	e.File("/index.html", "index.html")

	err := e.Start(":8000")
	if err != nil {
		panic(err)
	}
}

func main() {
	log.Println("service started")

	redis.SetupRedisClient()

	utils.StartRabbitMQ()

	go utils.RMQ.RecieveMessages()

	setupEndpoints()
}
