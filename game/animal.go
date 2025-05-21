package game

import (
	"math"
	"math/rand"
	"sync"

	_ "image/png"

	"example.com/artificial-life/neat"
	"github.com/gopxl/pixel/v2"
)

type TurningState uint8

const (
	LEFT TurningState = iota
	RIGHT
	STRAIGHT
)

type AnimalType uint8

const (
	PREY AnimalType = iota
	HUNTER
)

const (
	HUNGER_PERIOD = float64(1.0)
)

type Animal struct {
	x               float64
	y               float64
	w               float64
	h               float64
	dirTheta        float64
	speed           float64
	turningState    TurningState
	turningRate     float64
	ticksUntilHurt  int
	fovRays         int
	fov             float64
	viewingDistance int
	hp              int
	animalType      AnimalType
	// mu              sync.Mutex
	ticksToAppear   int
	fitnessGoal     int
	fitness         int
	reproCoolDown   int
	brain           *neat.Genome
	sprite          *pixel.Sprite
}

func InitAnimal(x, y float64, animalType AnimalType) *Animal {

	spriteFile := ""
	if animalType == PREY {
		spriteFile = "./img/looking-left.png"
	} else {
		spriteFile = "./img/fox.png"
	}

	pic, err := LoadPicture(spriteFile)
	if err != nil {
		panic(err)
	}

	sprite := pixel.NewSprite(pic, pic.Bounds())

	w := sprite.Frame().Max.X - sprite.Frame().Min.X
	h := sprite.Frame().Max.Y - sprite.Frame().Min.Y
	fov := math.Pi / 2
	fovRays := 6

	brain := neat.CreateGenome(0, fovRays+1, 3)
	brain.InitializeFromInitialConfig()

	return &Animal{x: x,
		y:               y,
		w:               w,
		h:               h,
		speed:           60,
		turningRate:     0.125,
		ticksUntilHurt:  20,
		hp:              20,
		fov:             fov,
		fovRays:         fovRays,
		viewingDistance: 50,
		animalType:      animalType,
		fitnessGoal:     Game.TicksPerSecond * 30,
		brain:           brain,
		sprite:          sprite,
	}
}

func (a *Animal) IncrementPos(dx, dy float64) {
	a.x += dx
	a.y += dy
}

func (a *Animal) GetPos() (x, y float64) {
	x = a.x
	y = a.y
	return x, y
}

func (a *Animal) GetDim() (w, h float64) {
	return a.w, a.h
}

func (a *Animal) GetTheta() float64 {
	return a.dirTheta
}

func (a *Animal) GetSpeed() float64 {
	return a.speed
}

func (a *Animal) GetSprite() *pixel.Sprite {
	return a.sprite
}

func (a *Animal) TurnDelta(dAngle float64) {
	a.dirTheta += dAngle
}

func (a *Animal) SetTurningState(x TurningState) {
	a.turningState = x
}

func (a *Animal) updateDirection() {
	a.see()
	a.brain.SetHpInput(a.hp)
	a.brain.Think()
	turnLeft := a.brain.GetOutput(0)
	turnRight := a.brain.GetOutput(1)
	noTurn := a.brain.GetOutput(2)
	// fmt.Println(
	// 	"ANDORRA",
	// 	a.brain.GetOutput(0),
	// 	a.brain.GetOutput(1),
	// 	a.brain.GetOutput(2),
	// )
	if turnLeft > turnRight && turnLeft > noTurn {
		a.turningState = LEFT
	}
	if turnRight > turnLeft && turnRight > noTurn {
		a.turningState = RIGHT
	}
	if noTurn > turnLeft && noTurn > turnRight {
		a.turningState = STRAIGHT
	}

	switch a.turningState {
	case LEFT:
		a.TurnDelta(a.turningRate)
	case RIGHT:
		a.TurnDelta(-a.turningRate)
	}
}

func (a *Animal) Update(wg *sync.WaitGroup) {
	defer wg.Done()
	a.ticksUntilHurt--
	if a.ticksUntilHurt <= 0 {
		a.hp--
		a.ticksUntilHurt = 20
	}

	a.fitness++

	if a.ticksUntilHurt > 20 * 10 && a.reproCoolDown <= 0 {
		newAnimal := a.Copy()
		newAnimal.TurnDelta(rand.NormFloat64() * math.Pi)
		newAnimal.hp = 20
		newAnimal.ticksUntilHurt = 100
		newAnimal.ticksToAppear = Game.TicksPerSecond + 1
		newAnimal.fitness = 0
		newAnimal.fitnessGoal = a.fitness
		newAnimal.brain.Mutate()

		newAnimals = append(newAnimals, newAnimal)


		if a.fitness > 20*40 {
			newAnimal := a.Copy()
			newAnimal.TurnDelta(rand.NormFloat64() * math.Pi)
			newAnimal.hp = 20
			newAnimal.ticksUntilHurt = 100
			newAnimal.ticksToAppear = Game.TicksPerSecond + 1
			newAnimal.fitness = 0
			newAnimal.fitnessGoal = a.fitness
			newAnimal.brain.Mutate()

			newAnimals = append(newAnimals, newAnimal)
		}
		a.reproCoolDown = 20*10
	} else {
		a.reproCoolDown--
	}


	if a.animalType == PREY {
		a.CheckNHandlePlantCollisions()
	}

	a.updateDirection()
	dx := math.Cos(a.dirTheta) * a.speed * float64(1.0/20.0)
	dy := math.Sin(a.dirTheta) * a.speed * float64(1.0/20.0)
	a.x = a.x + dx
	a.y = a.y + dy
}

func (a *Animal) see() {
	fov := a.fov
	rays := float64(a.fovRays)

	i := 0
	for angle := a.dirTheta - fov/2; angle <= a.dirTheta+fov/2; angle += fov / rays {
		a.castRay(i, angle)
		i++
	}
}

func (a *Animal) castRay(idx int, theta float64) {
	vec := pixel.ZV
	vec.X = a.x + a.w/2
	vec.Y = a.y + a.h/2

	dir := pixel.ZV
	dir.X = math.Cos(theta) * a.w / 2
	dir.Y = math.Sin(theta) * a.w / 2

	for i := 1; i <= a.viewingDistance; i++ {
		vec = vec.Add(dir)
		x := vec.X - math.Mod(vec.X, 16.0)
		y := vec.Y - math.Mod(vec.Y, 16.0)
		// pos := pixel.ZV
		// pos.X = x
		// pos.Y = y
		// Squares = append(Squares, pos)
		food := FoodBlocks[HashCoords(x, y)]
		if food != nil && food.fp > 0 {
			a.brain.SetVisionInput(idx, 1-float64(i)/float64(a.viewingDistance+1))
			return
		}
	}
	a.brain.SetVisionInput(idx, 1000_000_000)
}

func (a *Animal) GetHP() int {
	return a.hp
}

func (a *Animal) CheckNHandlePlantCollisions() {
	if IsCellFull(a.x, a.y) {
		sx := a.x - math.Mod(a.x, 16.0)
		sy := a.y - math.Mod(a.y, 16.0)
		food := FoodBlocks[HashCoords(sx, sy)]
		if food.Eeat() > 0 {
			a.ticksUntilHurt += 100
			a.hp++
		}
	} else if IsCellFull(a.x+16, a.y) {
		sx := a.x + 16 - math.Mod(a.x, 16.0)
		sy := a.y - math.Mod(a.y, 16.0)
		food := FoodBlocks[HashCoords(sx, sy)]
		if food.Eeat() > 0 {
			a.ticksUntilHurt += 100
			a.hp++
		}
	} else if IsCellFull(a.x, a.y+16) {
		sx := a.x - math.Mod(a.x, 16.0)
		sy := a.y + 16 - math.Mod(a.y, 16.0)
		food := FoodBlocks[HashCoords(sx, sy)]
		if food.Eeat() > 0 {
			a.ticksUntilHurt += 100
			a.hp++
		}
	} else if IsCellFull(a.x+16, a.y+16) {
		sx := a.x + 16 - math.Mod(a.x, 16.0)
		sy := a.y + 16 - math.Mod(a.y, 16.0)
		food := FoodBlocks[HashCoords(sx, sy)]
		if food.Eeat() > 0 {
			a.ticksUntilHurt += 100
			a.hp++
		}
	}
}

// func (a *Animal) CheckNHandlePreyCollisions() {
// 	for _, other := range Animals {
// 		x, y := other.GetPos()
// 		w, h := other.GetDim()
// 		if a.Collides(x, y, w, h) {
// 			a.hp += other.GetEaten()
// 			a.ticksUntilHurt += 100
// 		}
// 	}
// }

// func (a *Animal) GetEaten() int {
// 	a.mu.Lock()
// 	defer a.mu.Unlock()

// 	if a.hp > 0 {
// 		a.hp = 0
// 		return 1
// 	}

// 	return 0
// }

func (a *Animal) Collides(x2, y2, w2, h2 float64) bool {
	x1 := a.x
	y1 := a.y
	w1 := a.w
	h1 := a.h

	if x1+w1 > x2 && x1 < x2+w2 && y1+h1 > y2 && y1 < y2+h2 {
		return true
	}

	return false
}

func (a Animal) Copy() *Animal {
	return &a
}

func (a *Animal) GetNumberOfNeurons() int {
	return a.brain.GetNumberOfNeurons()
}
