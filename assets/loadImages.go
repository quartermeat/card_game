package assets

import (
	"encoding/csv"
	"image"
	"io"
	"os"
	"strconv"

	"github.com/faiface/pixel"
	"github.com/pkg/errors"
)

// ObjectImageAsset holds images for objects
type ObjectImageAsset struct {
	Description string
	Sheet       pixel.Picture
	Images      map[string]pixel.Rect
	ImageKeys   []string
}

func (objImage ObjectImageAsset) GetDescription() string {
	return objImage.Description
}

func (objImage ObjectImageAsset) GetSheet() pixel.Picture {
	return objImage.Sheet
}

func (objImage ObjectImageAsset) GetImages() map[string]pixel.Rect {
	return objImage.Images
}

func (m ObjectImageAsset) GetAnims() map[string][]pixel.Rect {
	panic("not implemented") // TODO: Implement
}

func (objImage ObjectImageAsset) GetKeys() []string {
	return objImage.ImageKeys
}

// getImages will take the sheet, a rectanlge for size, offset_x is the vertical space in-between images in pixels,
// offset_y is the vertical space in-between images in pixels
func getImages(sheet pixel.Picture, imageRect pixel.Rect, offset_x float64, offset_y float64) map[int][]pixel.Rect {
	images := make(map[int][]pixel.Rect, 0)
	row := 0
	for y := imageRect.Min.Y; y <= sheet.Bounds().Max.Y; y += imageRect.H() + offset_y {
		for x := imageRect.Min.X; x <= sheet.Bounds().Max.X; x += imageRect.W() + offset_x {
			images[row] = append(images[row], pixel.R(x, y, x+imageRect.W(), y+imageRect.H()))
		}
		row++
	}
	return images
}

// AddAssetsFromMultipleSheets adds assets
// imageRect is the size of the images being loaded as a rect
// offset_x is the horizontal space between images
// offset_y is the vertical space between images
func (objectAssets ObjectAssets) AddImageAssets(imageRect pixel.Rect, offset_x float64, offset_y float64, sheetDesc string, sheetPath, descPath string) (ObjectAssets, error) {
	var (
		err       error
		sheet     pixel.Picture
		images    map[string]pixel.Rect
		imageKeys []string
	)

	// total hack, nicely format the error at the end, so I don't have to type it every time
	defer func() {
		if err != nil {
			err = errors.Wrap(err, "error loading animation sheet")
		}
	}()

	// open and load the spritesheet
	sheetFile, err := os.Open(sheetPath)
	if err != nil {
		return nil, err
	}
	defer sheetFile.Close()
	sheetImg, _, err := image.Decode(sheetFile)
	if err != nil {
		return nil, err
	}
	sheet = pixel.PictureDataFromImage(sheetImg)

	// create a map of rects from sheet
	tempImages := getImages(sheet, imageRect, offset_x, offset_y)

	descFile, err := os.Open(descPath)
	if err != nil {
		return nil, err
	}
	defer descFile.Close()

	images = make(map[string]pixel.Rect)
	imageKeys = make([]string, 0)

	// load the animation information, name and interval inside the spritesheet
	desc := csv.NewReader(descFile)
	desc.Comma = ','
	desc.Comment = '#'
	for {
		image, err := desc.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		name := image[0]
		row, _ := strconv.Atoi(image[1])
		col, _ := strconv.Atoi(image[2])

		images[name] = tempImages[row][col]

		imageKeys = append(imageKeys, name)

	}
	newAsset := new(ObjectImageAsset)
	newAsset.Description = sheetDesc
	newAsset.ImageKeys = imageKeys
	newAsset.Sheet = sheet
	newAsset.Images = images

	objectAssets = append(objectAssets, *newAsset)

	return objectAssets, err
}
