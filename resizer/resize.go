package resizer

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/nfnt/resize"
)

func (r *Resize) ResizeImage(name string, qualityIndex float32) {
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
