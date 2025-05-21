package game

import (
	"image"
	"math"
	"math/rand"
	"os"

	_ "image/png"

	"github.com/gopxl/pixel/v2"
)

type GameConfig struct {
	TicksPerSecond int
}

var Game GameConfig
var Animals []*Animal
var FoodBlocks = make(map[int]*Food)
var Squares []pixel.Vec

func HashCoords(x, y float64) int {
	return int(x) * 100_000 + int(y)
}

func InitGameConfig() {
	Game = GameConfig{TicksPerSecond: 20}
	Animals = append(Animals, InitAnimal())

	for x := 0; x < 4096; x+=32 {
	    for y := 0; y < 4096; y+=32 {
		    if rand.Float64() < 0.05 {
			    FoodBlocks[HashCoords(float64(x), float64(y))] = InitFood(float64(x), float64(y))
		    }
	    }
	}
}

func LoadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func PruneDeadAnimals() {
	for k, animal := range Animals {
		if animal.GetHP() <= 0 {
			Animals = append(Animals[:k], Animals[k+1:]...)
		}
	}
}

func IsCellFull(x, y float64) bool {
	sx := x - math.Mod(x, 16.0)
	sy := y - math.Mod(y, 16.0)
	food := FoodBlocks[HashCoords(sx, sy)]

	return food != nil
}
