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
var newAnimals []*Animal
var FoodBlocks = make(map[int]*Food)
var FoodBlocksToGrow []*Food
var Squares []pixel.Vec

func HashCoords(x, y float64) int {
	return int(x)*100_000 + int(y)
}

func InitGameConfig() {
	Game = GameConfig{TicksPerSecond: 60}

	for x := 0; x < 4096; x += 32 {
		for y := 0; y < 4096; y += 32 {
			if rand.Float64() < 0.05 {
				FoodBlocks[HashCoords(float64(x), float64(y))] = InitFood(float64(x), float64(y))
			}
			if rand.Float64() < 0.005 {
				r := rand.Float64()
				var animalType AnimalType
				if r < 0.0 {
					animalType = HUNTER
				} else {
					animalType = PREY
				}
				Animals = append(Animals, InitAnimal(float64(x), float64(y), animalType))
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
	for k := 0; k < len(Animals); k++ {
		if Animals[k].GetHP() <= 0 {
			if len(Animals) < 5 {
				newAnimal := Animals[k].Copy()
				newAnimal.TurnDelta(rand.NormFloat64() * math.Pi)
				newAnimal.hp = 20
				newAnimal.ticksUntilHurt = 100
				newAnimal.ticksToAppear = Game.TicksPerSecond + 1
				newAnimal.fitness = 0
				newAnimal.fitnessGoal = Animals[k].fitness
				newAnimal.brain.MutateHighVariability()

				newAnimals = append(newAnimals, newAnimal)
			}
			Animals = append(Animals[:k], Animals[k+1:]...)
		}
	}

	for k := 0; k < len(newAnimals); k++ {
		animal := newAnimals[k]
		animal.ticksToAppear--
		if animal.ticksToAppear <= 0 {
			Animals = append(Animals, animal)
			newAnimals = append(newAnimals[:k], newAnimals[k+1:]...)
		}
	}
}

func RegrowPlants() {
	for i := 0; i < len(FoodBlocksToGrow); i++ {
		v := FoodBlocksToGrow[i]
		v.ticksToRegrow--
		if v.ticksToRegrow <= 0 {
			v.currSpriteIdx = 0
			v.fp = 1

			FoodBlocksToGrow = append(FoodBlocksToGrow[:i], FoodBlocksToGrow[i+1:]...)
		}
	}
}

func IsCellFull(x, y float64) bool {
	sx := x - math.Mod(x, 16.0)
	sy := y - math.Mod(y, 16.0)
	food := FoodBlocks[HashCoords(sx, sy)]

	return food != nil
}
