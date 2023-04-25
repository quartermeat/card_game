// The provided Go code is part of a card game application, specifically a package that defines the rules and assets for the game. Let me break down the key parts of the code:
//  1. The `import` statement imports required packages, including the `Pixel` library for 2D graphics and a custom `assets` package.
//  2. The `const` block defines several constants related to the game's appearance, such as the window size, card dimensions, gaps between cards, and paths to various asset // files, such as images and metadata.
//  3. `WINDOW_SIZE` is a `pixel.Rect` variable that defines the window size for the game using the `window_width` and `window_height` constants.
//  4. `CardRect` is a `pixel.Rect` variable that defines the dimensions of a card using the `MinXCoord`, `MinYCoord`, `MaxXCoord`, and `MaxYCoord` constants.
//  5. The `LoadAssets` function loads the game's assets, including images and animations, into an `assets.ObjectAssets` instance. It takes a `sysErrors` parameter, which is a slice of errors to keep track of any issues that occur during asset loading. The function returns an updated `assets.ObjectAssets` instance with the loaded assets.
//     The function iterates through a series of asset descriptions and paths, calling the appropriate methods (`AddAnimationAssets` and `AddImageAssets`) to load each asset. If any errors occur during this process, they are appended to the `sysErrors` slice.
//     At the end of the function, it checks the `sysErrors` slice for any errors. If any are present, the function panics and prints the error, halting the program. Otherwise, it returns the `objectAssets` instance with the loaded assets.
//
// This package is responsible for defining the game's appearance, managing assets, and providing an interface to load and access these assets during the game's runtime.
package card_game_rules

import (
	"github.com/faiface/pixel"
	"github.com/quartermeat/card_game/assets"
)

// domain specific constants
const (
	APP_TITLE     string  = "Card Game"
	window_height float64 = 768
	window_width  float64 = 1366
	// coordinates of the top left image
	MinXCoord float64 = 26
	MinYCoord float64 = 71
	MaxXCoord float64 = 250
	MaxYCoord float64 = 395
	// gaps that exist between images
	HorizontalGap float64 = 26
	VerticalGap   float64 = 26
	//cursor
	CURSOR_ANIMATIONS_DESC = "hand"
	CURSOR_SPRITE_SHEET    = "assets/animations/cursor/cursorHandBig.png"
	CURSOR_ICON_SIZE       = 32
	CURSOR_META            = "assets/animations/cursor/cursorAnimations.csv"
	// background
	BACKGROUND_IMAGE = "assets/images/background/woodenTable.png"
	// actions3
	HAM_RADIO_DESC  = "ham_radio"
	HAM_RADIO_IMAGE = "assets/images/zombieCards/1xHamRadio.png"
	HAM_RADIO_META  = "assets/images/zombieCards/hamRadio.csv"
	TRASH_DESC      = "trash"
	TRASH_IMAGE     = "assets/images/zombieCards/1xTrash.png"
	TRASH_META      = "assets/images/zombieCards/trash.csv"
	SLUGS_DESC      = "slugs"
	SLUGS_IMAGE     = "assets/images/zombieCards/6xSlugs.png"
	SLUGS_META      = "assets/images/zombieCards/slugs.csv"
	ZOMBIES_DESC    = "zombies"
	ZOMBIES_IMAGE   = "assets/images/zombieCards/6xZombies.png"
	ZOMBIES_META    = "assets/images/zombieCards/zombies.csv"
	ACTIONS_1_DESC  = "actions1"
	ACTIONS_1_IMAGE = "assets/images/zombieCards/10xActions1.png"
	ACTIONS_1_META  = "assets/images/zombieCards/actions1.csv"
	ACTIONS_2_DESC  = "actions2"
	ACTIONS_2_IMAGE = "assets/images/zombieCards/10xActions2.png"
	ACTIONS_2_META  = "assets/images/zombieCards/actions2.csv"
	ACTIONS_3_DESC  = "actions3"
	ACTIONS_3_IMAGE = "assets/images/zombieCards/10xActions3.png"
	ACTIONS_3_META  = "assets/images/zombieCards/actions3.csv"
	BULLETS_DESC    = "bullets"
	BULLETS_IMAGE   = "assets/images/zombieCards/15xBullets.png"
	BULLETS_META    = "assets/images/zombieCards/bullets.csv"
)

// WINDOW_SIZE is the size of the window to be created by pixel gl
var WINDOW_SIZE pixel.Rect = pixel.R(0, 0, window_width, window_height)

// top left card in image
var CardRect pixel.Rect = pixel.R(MinXCoord, MinYCoord, MaxXCoord, MaxYCoord)

// LoadAssets loads assets used for the application
func LoadAssets(sysErrors []error) assets.ObjectAssets {
	var objectAssets assets.ObjectAssets
	objectAssets, err1 := objectAssets.AddAnimationAssets(CURSOR_ANIMATIONS_DESC, CURSOR_SPRITE_SHEET, CURSOR_META, CURSOR_ICON_SIZE)
	objectAssets, err2 := objectAssets.AddImageAssets(CardRect, HorizontalGap, VerticalGap, ACTIONS_1_DESC, ACTIONS_1_IMAGE, ACTIONS_1_META)
	objectAssets, err3 := objectAssets.AddImageAssets(CardRect, HorizontalGap, VerticalGap, ACTIONS_2_DESC, ACTIONS_2_IMAGE, ACTIONS_2_META)
	objectAssets, err4 := objectAssets.AddImageAssets(CardRect, HorizontalGap, VerticalGap, ACTIONS_3_DESC, ACTIONS_3_IMAGE, ACTIONS_3_META)
	objectAssets, err5 := objectAssets.AddImageAssets(CardRect, HorizontalGap, VerticalGap, HAM_RADIO_DESC, HAM_RADIO_IMAGE, HAM_RADIO_META)
	objectAssets, err6 := objectAssets.AddImageAssets(CardRect, HorizontalGap, VerticalGap, TRASH_DESC, TRASH_IMAGE, TRASH_META)
	objectAssets, err7 := objectAssets.AddImageAssets(CardRect, HorizontalGap, VerticalGap, SLUGS_DESC, SLUGS_IMAGE, SLUGS_META)
	objectAssets, err8 := objectAssets.AddImageAssets(CardRect, HorizontalGap, VerticalGap, ZOMBIES_DESC, ZOMBIES_IMAGE, ZOMBIES_META)
	objectAssets, err9 := objectAssets.AddImageAssets(CardRect, HorizontalGap, VerticalGap, BULLETS_DESC, BULLETS_IMAGE, BULLETS_META)
	sysErrors = append(sysErrors, err1)
	sysErrors = append(sysErrors, err2)
	sysErrors = append(sysErrors, err3)
	sysErrors = append(sysErrors, err4)
	sysErrors = append(sysErrors, err5)
	sysErrors = append(sysErrors, err6)
	sysErrors = append(sysErrors, err7)
	sysErrors = append(sysErrors, err8)
	sysErrors = append(sysErrors, err9)
	for _, sysError := range sysErrors {
		if sysError != nil {
			panic(sysError)
		}
	}
	return objectAssets
}
