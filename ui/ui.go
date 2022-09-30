package ui

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/golang/freetype/truetype"
	"github.com/quartermeat/card_game/input"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
)

const (
	TRUETYPE_FONT_PATH = "assets\\fonts\\intuitive.ttf"
)

type GUI struct {
	atlas    *text.Atlas
	txt      *text.Text
	face     font.Face
	doUpdate bool
	lines    []string
}

type IGUI interface {
	InitGUI()
	UpdateGUI()
	DrawGUI()
}

func (gui *GUI) loadTTF(path string, size float64) (font.Face, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	font, err := truetype.Parse(bytes)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(font, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	}), nil
}

func (gui *GUI) InitGUI() {

	gui.face, _ = gui.loadTTF(TRUETYPE_FONT_PATH, 52)
	gui.atlas = text.NewAtlas(gui.face, text.ASCII)
	gui.txt = text.New(pixel.V(0, 0), gui.atlas)
	gui.doUpdate = true
	gui.txt.Color = colornames.Green
	gui.lines = []string{}
}

// remove index s from slice, return new slice
func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

// UpdateGUI does gui updates based on game commands
func (gui *GUI) UpdateGUI(cmds input.Commands) {

	for key, _ := range cmds {
		gui.lines = append(gui.lines, fmt.Sprintf("executing: %s", key))
	}

	for idx, line := range gui.lines {
		fmt.Fprintln(gui.txt, line)
		gui.lines = remove(gui.lines, idx)
	}
}

// DrawGUI draws the gui on the specified window
func (gui *GUI) DrawGUI(win *pixelgl.Window, cam *pixel.Matrix) {

	gui.txt.Draw(win, pixel.IM.Moved(cam.Unproject(win.Bounds().Min).Sub(gui.txt.Bounds().Min)))
}
