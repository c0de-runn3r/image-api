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
	"path"
	"path/filepath"
	"time"

	"github.com/labstack/echo"
	"github.com/nfnt/resize"
)

var (
	StoragePath    string
	AllowedFileExt = []string{".jpg", ".jpeg", ".png", ".gif"}
)

type Image struct {
	ID   int
	Name string
}

func UploadImage(c echo.Context) error {
	log.Println("got new image upload request")

	ogrFilename, newFilename := copyFile(c)

	log.Println("file uploaded successfully")

	idStr := generateID()

	redis.Redis.AddPair(idStr, newFilename)
	RMQ.SendMessage(idStr)

	return c.HTML(http.StatusOK, fmt.Sprintf("<p>Image %s uploaded successfully. ID: %s</p>", ogrFilename, idStr))
}

func SendImage(c echo.Context) error {
	id := c.FormValue("id")
	quality := c.FormValue("quality")
	fileName := redis.Redis.GetValue(id)
	switch quality {
	case "100":
		return c.File(filepath.Join(StoragePath, fileName))
	case "75":
		return c.File(filepath.Join(StoragePath, "0.75_"+fileName))
	case "50":
		return c.File(filepath.Join(StoragePath, "0.5_"+fileName))
	case "25":
		return c.File(filepath.Join(StoragePath, "0.25_"+fileName))
	default:
		return c.File(filepath.Join(StoragePath, fileName))
	}
}
func ResizeImage(name string, qualityIndex float32) {
	path := path.Join(StoragePath, name)
	file, err := os.Open(path)
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

	buffer := resize.Resize(width, height, img, resize.Lanczos3)

	newName := fmt.Sprintf("%v_%s", qualityIndex, name)
	newPath := filepath.Join(StoragePath, newName)
	out, err := os.Create(newPath)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	switch format {
	case "jpg", "jpeg":
		if err := jpeg.Encode(out, buffer, nil); err != nil {
			log.Panic(err)
		}
	case "png":
		if err := png.Encode(out, buffer); err != nil {
			log.Panic(err)
		}
	case "gif":
		if err := gif.Encode(out, buffer, nil); err != nil {
			log.Panic(err)
		}
	default:
		if err := jpeg.Encode(out, buffer, nil); err != nil {
			log.Panic(err)
		}
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

func copyFile(c echo.Context) (string, string) {
	_, err := os.ReadDir(StoragePath)
	if err != nil {
		if err := os.Mkdir(StoragePath, 0755); err != nil {
			log.Panicln("error occured creating directory", err)
		}
	}

	file, err := c.FormFile("image")
	if err != nil {
		log.Printf("error occured uploading file: %e", err)
	}

	if !checkExtention(file.Filename) {
		c.HTML(http.StatusUnsupportedMediaType, "<p>Only JPG, JPEG, PNG and GIF formats are allowed!")
	}

	src, err := file.Open()
	if err != nil {
		log.Println("error reading from file", err)
	}
	defer src.Close()

	formattedName := fmt.Sprintf("%x_%s", time.Now().UnixMilli(), file.Filename) // to prevent conflicts when two files have the same names
	path := filepath.Join(StoragePath, formattedName)
	dst, err := os.Create(path)
	if err != nil {
		log.Println("error creating file", err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		log.Println("error copying to file", err)
	}
	return file.Filename, formattedName
}

func originalSize(file image.Image) (uint, uint) {
	bounds := file.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	return uint(width), uint(height)
}

func checkExtention(fileName string) bool {
	fileExtension := filepath.Ext(fileName)
	for _, v := range AllowedFileExt {
		if fileExtension == v {
			return true
		}
	}
	return false
}

func generateID() string {
	rand.Seed(time.Now().UnixNano())
	id := rand.Intn(999_999-100_000) + 100_000
	return fmt.Sprintf("%v", id)
}
