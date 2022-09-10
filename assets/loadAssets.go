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

// type ObjectAsset struct {
// 	Sheet    pixel.Picture
// 	Anims    map[string][]pixel.Rect
// 	AnimKeys []string
// }

const MouseIconPixelSize float64 = 16
const CardImageSize float64 = 368
const (
	CursorAnimations string = "hand"
	TestCard         string = "test_card"
)

//ObjectAssets holds images for objects
type ObjectAsset struct {
	Description string
	Sheet       pixel.Picture
	Anims       map[string][]pixel.Rect
	AnimKeys    []string
}

type ObjectAssets []ObjectAsset

func getFrames(sheet pixel.Picture, frameSize float64) [][]pixel.Rect {
	frames := make([][]pixel.Rect, 0)
	for y := 0.0; y+frameSize <= sheet.Bounds().Max.Y; y += frameSize {
		temp := make([]pixel.Rect, 0)
		for x := 0.0; x+frameSize <= sheet.Bounds().Max.X; x += frameSize {
			temp = append(temp, pixel.R(x, y, x+frameSize, y+frameSize))
		}
		frames = append(frames, temp)
	}
	return frames
}

func (objectAssets ObjectAssets) AddAssets(sheetDesc string, sheetPath, descPath string, frameSize float64) (ObjectAssets, error) {
	var (
		err      error
		sheet    pixel.Picture
		anims    map[string][]pixel.Rect
		animKeys []string
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

	// create a slice of frames inside the spritesheet
	frames := getFrames(sheet, frameSize)

	descFile, err := os.Open(descPath)
	if err != nil {
		return nil, err
	}
	defer descFile.Close()

	anims = make(map[string][]pixel.Rect)
	animKeys = make([]string, 0)

	// load the animation information, name and interval inside the spritesheet
	desc := csv.NewReader(descFile)
	desc.Comma = ','
	desc.Comment = '#'
	for {
		anim, err := desc.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		name := anim[0]
		row, _ := strconv.Atoi(anim[1])
		start, _ := strconv.Atoi(anim[2])
		end, _ := strconv.Atoi(anim[3])

		anims[name] = frames[row][start : end+1]

		animKeys = append(animKeys, name)

	}
	newAsset := new(ObjectAsset)
	newAsset.Description = sheetDesc
	newAsset.AnimKeys = animKeys
	newAsset.Sheet = sheet
	newAsset.Anims = anims

	objectAssets = append(objectAssets, *newAsset)

	return objectAssets, err
}
