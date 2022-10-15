package camera

import (
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// TODO: camera needs to be fixed
type Camera struct {
	Matrix    *pixel.Matrix
	Position  pixel.Vec
	Speed     float64
	Zoom      float64
	ZoomSpeed float64
}

type ICamera interface {
	SetupCamera(win *pixelgl.Window)
	GetMatrix() *pixel.Matrix
	GetSpeed() float64
	GetZoom() float64
	SetZoom(zoom float64)
	GetZoomSpeed() float64
	GetPosition() *pixel.Vec
}

func (cam *Camera) SetZoom(win *pixelgl.Window) {
	cam.Zoom *= math.Pow(cam.ZoomSpeed, win.MouseScroll().Y)
}

// SetupCamera sets up the game camera using the current window
func (cam *Camera) SetupCamera(win *pixelgl.Window) {
	cam.Position = pixel.ZV
	cam.Speed = 500.0
	cam.Zoom = 1.0
	cam.ZoomSpeed = 1.2
	// cam.Matrix = pixel.IM.Scaled(cam.Position, cam.Zoom).Moved(win.Bounds().Center().Sub(cam.Position))
}

func (cam *Camera) GetSpeed() *float64 {
	return &(cam.Speed)
}

func (cam *Camera) GetZoom() *float64 {
	return &(cam.Zoom)
}

func (cam *Camera) GetZoomSpeed() *float64 {
	return &cam.ZoomSpeed
}
