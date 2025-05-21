package game

import (
	"sync"

	_ "image/png"

	"github.com/gopxl/pixel/v2"
)

type Food struct {
	x             float64
	y             float64
	w             float64
	h             float64
	fp            int
	maxFp         int
	sprites       []*pixel.Sprite
	currSpriteIdx int
	mu            sync.Mutex
	available     bool
}

func InitFood(x, y float64) *Food {
	w := 16.0
	h := 16.0

	spritesheet, err := LoadPicture("./img/shrubs.png")
	if err != nil {
		panic(err)
	}

	sprites := make([]*pixel.Sprite, 2)
	var frames []pixel.Rect
	for x := spritesheet.Bounds().Min.X; x < spritesheet.Bounds().Max.X; x += w {
		for y := spritesheet.Bounds().Min.Y; y < spritesheet.Bounds().Max.Y; y += h {
			frames = append(frames, pixel.R(x, y, x+w, y+h))
		}
	}

	sprites[0] = pixel.NewSprite(spritesheet, frames[0])
	sprites[1] = pixel.NewSprite(spritesheet, frames[1])
	return &Food{x: x, y: y, w: w, h: h, fp: 1, maxFp: 10, sprites: sprites}
}

func (f *Food) GetPos() (x, y float64) {
	return f.x, f.y
}

func (f *Food) GetDim() (w, h float64) {
	return f.w, f.h
}

func (f *Food) GetCurrSprite() *pixel.Sprite {
	return f.sprites[f.currSpriteIdx]
}

func (f *Food) Eeat() int {
	if f == nil {
		return 0
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	fpEaten := 0
	if f.fp > 0 {
		f.fp--
		if f.fp == 0 {
			f.currSpriteIdx = 1
		}
		fpEaten = 1
	}


	return fpEaten
}
