package assets

import (
	"image"
	"os"

	"github.com/faiface/pixel"
)

// LoadPicture returns a pixel.Picture from a path
func LoadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return pixel.PictureDataFromImage(img), nil
}
