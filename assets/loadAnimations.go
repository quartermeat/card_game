// Package 'assets' is used to load sprite animations
package assets

import (
	"encoding/csv"
	"fmt"
	"image"
	"io"
	"os"
	"strconv"

	"github.com/faiface/pixel"
	"github.com/pkg/errors"
)

// ObjectAnimationAsset holds animations for objects
type ObjectAnimationAsset struct {
	Description string
	Sheet       pixel.Picture
	Anims       map[string][]pixel.Rect
	AnimKeys    []string
}

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

func (animation ObjectAnimationAsset) GetDescription() string {
	return animation.Description
}

func (animation ObjectAnimationAsset) GetSheet() pixel.Picture {
	return animation.Sheet
}

func (animation ObjectAnimationAsset) GetImages() map[string]pixel.Rect {
	fmt.Printf("GetImages() is not implemented for ObjectAnimationAsset:%s", animation.Description)
	return nil
}

func (animation ObjectAnimationAsset) GetAnims() map[string][]pixel.Rect {
	return animation.Anims
}

func (animation ObjectAnimationAsset) GetKeys() []string {
	return animation.AnimKeys
}

func (objectAssets ObjectAssets) AddAnimationAssets(sheetDesc string, sheetPath, descPath string, frameSize float64) (ObjectAssets, error) {
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
	newAsset := new(ObjectAnimationAsset)
	newAsset.Description = sheetDesc
	newAsset.AnimKeys = animKeys
	newAsset.Sheet = sheet
	newAsset.Anims = anims

	objectAssets = append(objectAssets, *newAsset)

	return objectAssets, err
}
