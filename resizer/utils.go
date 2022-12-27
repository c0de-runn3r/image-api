package resizer

import (
	"fmt"
	"image"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo"
)

var (
	StoragePath    string
	AllowedFileExt = []string{".jpg", ".jpeg", ".png", ".gif"}
)

type Image struct {
	ID   int
	Name string
}

func UploadImage(c echo.Context) (string, string, string) {
	orgFilename, newFilename := copyFile(c)
	log.Println("file uploaded successfully")

	idStr := generateID()

	return idStr, orgFilename, newFilename
}

func SendImage(quality string, fileName string) string {
	switch quality {
	case "100":
		return filepath.Join(StoragePath, fileName)
	case "75":
		return filepath.Join(StoragePath, "0.75_"+fileName)
	case "50":
		return filepath.Join(StoragePath, "0.5_"+fileName)
	case "25":
		return filepath.Join(StoragePath, "0.25_"+fileName)
	default:
		return filepath.Join(StoragePath, fileName)
	}
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
		err := c.HTML(http.StatusUnsupportedMediaType, "<p>Only JPG, JPEG, PNG and GIF formats are allowed!")
		if err != nil {
			log.Println(err)
		}
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
