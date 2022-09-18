package app

import (
	"github.com/faiface/pixel"
	"github.com/quartermeat/card_game/assets"
)

// domain specific constants
const (
	APP_TITLE              string  = "Card Game"
	window_height          float64 = 1080
	window_width           float64 = 1290
	MinXCoord              float64 = 26
	MinYCoord              float64 = 71
	MaxXCoord              float64 = 250
	MaxYCoord              float64 = 395
	HorizontalGap          float64 = 26
	VerticalGap            float64 = 26
	CURSOR_ANIMATIONS_DESC         = "hand"
	CURSOR_SPRITE_SHEET            = "assets/animations/cursor/cursorHand.png"
	CURSOR_ICON_SIZE               = 16
	CURSOR_META                    = "assets/animations/cursor/cursorAnimations.csv"
	ACTIONS_3_DESCRIPTION          = "actions3"
	ACTIONS_3_IMAGE                = "assets/images/zombieCards/10xActions3.png"
	ACTIONS_3_META                 = "assets/images/zombieCards/actions3.csv"
)

// WINDOW_SIZE is the size of the window to be created by pixel gl
var WINDOW_SIZE pixel.Rect = pixel.R(0, 0, window_width, window_height)

// top left card in image
var CardRect pixel.Rect = pixel.R(MinXCoord, MinYCoord, MaxXCoord, MaxYCoord)

func loadAssets(sysErrors []error) assets.ObjectAssets {
	var objectAssets assets.ObjectAssets
	objectAssets, err1 := objectAssets.AddAnimationAssets(CURSOR_ANIMATIONS_DESC, CURSOR_SPRITE_SHEET, CURSOR_META, CURSOR_ICON_SIZE)
	objectAssets, err2 := objectAssets.AddImageAssets(CardRect, HorizontalGap, VerticalGap, ACTIONS_3_DESCRIPTION, ACTIONS_3_IMAGE, ACTIONS_3_META)
	sysErrors = append(sysErrors, err1)
	sysErrors = append(sysErrors, err2)
	for _, sysError := range sysErrors {
		if sysError != nil {
			panic(sysError)
		}
	}
	return objectAssets
}
