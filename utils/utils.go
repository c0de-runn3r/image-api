package utils

import (
	"fmt"
	"image"
	"image-api/redis"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/nfnt/resize"
)

const (
	initDirName = "images"
)

type Image struct {
	ID   int
	Name string
}

func UploadImage(wr http.ResponseWriter, req *http.Request) {
	log.Println("got new image upload request")
	// req.ParseMultipartForm(32 << 20) // TODO can limit memory of file

	_, err := os.ReadDir(initDirName)
	if err != nil {
		os.Mkdir(initDirName, 0755)
	}

	file, handler, err := req.FormFile("image")
	if err != nil {
		log.Printf("error occured uploading file: %e", err)
	}
	defer file.Close()

	newFile, err := os.CreateTemp(initDirName, "*-"+handler.Filename)
	if err != nil {
		log.Println(err)
	}
	defer newFile.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		log.Println(err)
	}
	newFile.Write(fileData)

	log.Println("file uploaded successfully")
	rand.Seed(time.Now().UnixNano())
	id := rand.Intn(999_999-100_000) + 100_000
	idStr := fmt.Sprintf("%v", id)
	redis.Redis.AddPair(idStr, newFile.Name())
	RMQ.SendMessage(idStr) // TODO message
}

func ResizeImage(nameWithPath string, qualityIndex float32) {
	// os.Chdir("..")
	file, err := os.Open(nameWithPath)
	if err != nil {
		log.Fatal(err)
	}
	img, format, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	origWidth, origHeight := originalSize(img)
	width := uint(float32(origWidth) * qualityIndex)
	height := uint(float32(origHeight) * qualityIndex)
	name := filepath.Base(nameWithPath)
	buffer := resize.Resize(width, height, img, resize.Lanczos3)
	newName := fmt.Sprintf("%v_%s", qualityIndex, name)
	newPath := filepath.Join(initDirName, newName)
	out, err := os.Create(newPath)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	switch format {
	case "jpg", "jpeg":
		jpeg.Encode(out, buffer, nil)
	case "png":
		png.Encode(out, buffer)
	case "gif":
		gif.Encode(out, buffer, nil)
	default:
		jpeg.Encode(out, buffer, nil)
	}
	log.Println("file resized successfully")
}

func OptimizeImages(id string) {
	log.Printf("optimizing image [id:%s]", id)
	file := redis.Redis.GetValue(id)
	ResizeImage(file, 0.75)
	ResizeImage(file, 0.50)
	ResizeImage(file, 0.25)
}

func originalSize(file image.Image) (uint, uint) {
	bounds := file.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	return uint(width), uint(height)
}
