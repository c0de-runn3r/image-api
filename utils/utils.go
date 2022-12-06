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

const (
	initDirName = "images"
)

type Image struct {
	ID   int
	Name string
}

func UploadImage(c echo.Context) error {
	log.Println("got new image upload request")

	_, err := os.ReadDir(initDirName)
	if err != nil {
		if err := os.Mkdir(initDirName, 0755); err != nil {
			panic(err)
		}
	}

	file, err := c.FormFile("image")
	if err != nil {
		log.Printf("error occured uploading file: %e", err)
	}

	fileExtension := filepath.Ext(file.Filename)
	if fileExtension == ".jpg" || fileExtension == ".jpeg" || fileExtension == ".png" || fileExtension == ".gif" {

		src, err := file.Open()
		if err != nil {
			log.Println(err)
		}
		defer src.Close()

		formattedName := fmt.Sprintf("%x_%s", time.Now().UnixMilli(), file.Filename)
		path := filepath.Join(initDirName, formattedName)
		dst, err := os.Create(path)
		if err != nil {
			log.Println(err)
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			log.Println(err)
		}

		log.Println("file uploaded successfully")

		rand.Seed(time.Now().UnixNano())
		id := rand.Intn(999_999-100_000) + 100_000
		idStr := fmt.Sprintf("%v", id)

		redis.Redis.AddPair(idStr, formattedName)
		RMQ.SendMessage(idStr) // TODO message
		return c.HTML(http.StatusOK, fmt.Sprintf("<p>Image %s uploaded successfully. ID: %s</p>", file.Filename, idStr))
	} else {
		return c.HTML(http.StatusUnsupportedMediaType, "<p>Only JPG, JPEG, PNG and GIF formats are allowed!")
	}

}

func SendImage(c echo.Context) error {
	id := c.FormValue("id")
	quality := c.FormValue("quality")
	file := redis.Redis.GetValue(id)
	switch quality {
	case "100":
		return c.File(filepath.Join(initDirName, file))
	case "75":
		return c.File(filepath.Join(initDirName, "0.75_"+file))
	case "50":
		return c.File(filepath.Join(initDirName, "0.5_"+file))
	case "25":
		return c.File(filepath.Join(initDirName, "0.25_"+file))
	default:
		return c.File(filepath.Join(initDirName, file))
	}
}
func ResizeImage(name string, qualityIndex float32) {
	// os.Chdir("..")
	path := path.Join(initDirName, name)
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
	newPath := filepath.Join(initDirName, newName)
	out, err := os.Create(newPath)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	switch format {
	case "jpg", "jpeg":
		if err := jpeg.Encode(out, buffer, nil); err != nil {
			panic(err)
		}
	case "png":
		if err := png.Encode(out, buffer); err != nil {
			panic(err)
		}
	case "gif":
		if err := gif.Encode(out, buffer, nil); err != nil {
			panic(err)
		}
	default:
		if err := jpeg.Encode(out, buffer, nil); err != nil {
			panic(err)
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

func originalSize(file image.Image) (uint, uint) {
	bounds := file.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	return uint(width), uint(height)
}
