package assets

import (
	"github.com/gopxl/pixel"
)

type IObjectAsset interface {
	GetDescription() string
	GetSheet() pixel.Picture
	GetImages() map[string]pixel.Rect
	GetAnims() map[string][]pixel.Rect
	GetKeys() []string
}

type ObjectAssets []IObjectAsset

// GetImage returns the image for the description
func (objectAssets ObjectAssets) GetImage(desc string) IObjectAsset {
	for _, objectAsset := range objectAssets {
		if _, ok := objectAsset.GetImages()[desc]; ok {
			return objectAsset
		}
	}
	return nil
}

func (objectAssets ObjectAssets) IsDescriptionAvailable(desc string) bool {
	for _, assets := range objectAssets {
		if desc == assets.GetDescription() {
			return true
		}
	}
	return false
}
