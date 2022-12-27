package server

import (
	"fmt"
	"log"
	"net/http"

	"image-api/queue"
	"image-api/resizer"
	"image-api/storage"

	"github.com/labstack/echo"
)

var (
	s *storage.StorageClient
	r resizer.Resize
	q *queue.QueueClient
)

func InitializeModules(storage *storage.StorageClient, resizer resizer.Resize, queue *queue.QueueClient) {
	s = storage
	r = resizer
	q = queue
}

func HandleUploadImage(c echo.Context) error {
	log.Println("got new image upload request")
	idStr, orgFilename, newFilename := resizer.UploadImage(c)
	s.AddPair(idStr, newFilename)
	q.SendID(idStr)

	return c.HTML(http.StatusOK, fmt.Sprintf("<p>Image %s uploaded successfully. ID: %s</p>", orgFilename, idStr))
}

func HandleSendImage(c echo.Context) error {
	id := c.FormValue("id")
	quality := c.FormValue("quality")
	fileName := s.GetValue(id)

	filepath := resizer.SendImage(quality, fileName)

	return c.File(filepath)
}

func HandleMessages() {
	for {
		ok, idStr := q.RecieveID()
		if ok {
			HandleOptimizeImages(idStr, &r, s)
		}
	}
}

func HandleOptimizeImages(idStr string, r *resizer.Resize, storage *storage.StorageClient) {
	log.Printf("optimizing image [id:%s]", idStr)
	file := storage.GetValue(idStr)
	OptimizeImages(file)
}

func OptimizeImages(file string) {
	r.ResizeImage(file, 0.75)
	r.ResizeImage(file, 0.50)
	r.ResizeImage(file, 0.25)
}
