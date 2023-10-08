package border

import (
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/imdraw"
	"github.com/gopxl/pixel/pixelgl"
)

type Color struct {
	R, G, B, A float64
}

// DrawBorder draws a border around the specified region with the specified color and width.
func DrawBorder(win *pixelgl.Window, region pixel.Rect, color pixel.RGBA, width float64) {
	// Calculate the four corners of the border.
	x1, y1 := region.Min.X, region.Min.Y
	x2, y2 := region.Max.X, region.Max.Y
	tl := pixel.V(x1, y2)
	tr := pixel.V(x2, y2)
	bl := pixel.V(x1, y1)
	br := pixel.V(x2, y1)

	// Draw the border using an imdraw.Drawer.
	imd := imdraw.New(nil)
	imd.Color = color
	imd.Push(tl, tr)
	imd.Line(width)
	imd.Push(tr, br)
	imd.Line(width)
	imd.Push(br, bl)
	imd.Line(width)
	imd.Push(bl, tl)
	imd.Line(width)
	imd.Draw(win)
}

// Example usage:
// func main() {
//     cfg := pixelgl.WindowConfig{
//         Title:  "My Border Test",
//         Bounds: pixel.R(0, 0, 800, 600),
//     }
//     win, err := pixelgl.NewWindow(cfg)
//     if err != nil {
//         panic(err)
//     }
//     for !win.Closed() {
//         win.Clear(colornames.Black)
//         DrawBorder(win, pixel.R(100, 100, 500, 400), colornames.White, 3)
//         win.Update()
//     }
// }
