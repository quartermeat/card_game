package card

import (
	"sync"

	"github.com/gopxl/pixel/pixelgl"
)

const (
	//card types
	HAM_RADIO         = "ham_radio"
	TRASH             = "trash"
	SLUG              = "slug"
	SURVIVOR          = "survivor"
	ZOMBIES           = "zombies"
	MORE_ZOMBIES      = "more_zombies"
	EVEN_MORE_ZOMBIES = "even_more_zombies"
	AMMO_BOX          = "ammo_box"
	BARRICADE         = "barricade"
	BULLET            = "bullet"
	COURAGE           = "courage"
	CUNNING           = "cunning"
	DECOY             = "decoy"
	HIDE              = "hide"
	HIGHER_GROUND     = "higher_ground"
	HOLLOW_POINTS     = "hollow_points"
	MAVERICK          = "maverick"
	MOLOTOV_COCKTAIL  = "molotov_cocktail"
	QUICK_ESCAPE      = "quick_escape"
	RECON             = "recon"
	REGROUP           = "regroup"
	RELOAD            = "reload"
	RESTOCK           = "restock"
	SACRIFICE         = "sacrifice"
	SCAVENGER         = "scavenger"
	SHOTGUN           = "shotgun"
	SIDEKICK          = "sidekick"
	STICK_TOGETHER    = "stick_together"
	ZOMBIE_SWARM      = "zombie_swarm"
	TACTICS           = "tactics"
	WEAPONS_CACHE     = "weapons_cache"
	CARD_BACK         = "card_back"
)


type ICard interface {
	Draw(*pixelgl.Window, bool, *sync.WaitGroup)
}

type IDeck interface {	
}