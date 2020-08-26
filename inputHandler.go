package main

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func removeObject(win *pixelgl.Window, cam *pixel.Matrix, gameObjs GameObjects, livingObjs LivingObjects, gibletObjs GibletObjects) {
	//handle removing an object
	if win.JustPressed(pixelgl.MouseButtonRight) {
		mouse := cam.Unproject(win.MousePosition())
		selectedObj, index, hit, err := gameObjs.getSelectedGameObj(mouse)
		if err != nil {
			fmt.Printf(err.Error())
		}
		if hit {
			fmt.Println("object id:", selectedObj.getID(), " removed")
			gameObjs = gameObjs.fastRemoveIndex(index)

			switch selectedObj.(type) {
			case *livingObject:
				{
					livingObjs = livingObjs.fastRemoveIndexFromLivingObjects(index)
					selectedObj = nil
				}
			case *gibletObject:
				{
					gibletObjs = gibletObjs.fastRemoveIndexFromGibletObjects(index)
					selectedObj = nil
				}
			}
		} else {
			fmt.Println("no object selected")
		}
		hit = false
	}
}
