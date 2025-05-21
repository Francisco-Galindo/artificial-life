package game

import (
	"math"

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
	brain           *neat.Genome
	sprite          *pixel.Sprite
}

func InitAnimal() *Animal {

	pic, err := LoadPicture("./img/looking-left.png")
	if err != nil {
		panic(err)
	}

	sprite := pixel.NewSprite(pic, pic.Bounds())

	w := sprite.Frame().Max.X - sprite.Frame().Min.X
	h := sprite.Frame().Max.Y - sprite.Frame().Min.Y
	fov := math.Pi / 2
	fovRays := 8

	brain := neat.CreateGenome(0, 8, 3)
	brain.InitializeFromInitialConfig()

	return &Animal{x: 2048,
		y:               2048,
		w:               w,
		h:               h,
		speed:           60,
		turningRate:     0.125,
		ticksUntilHurt:  20,
		hp:              20,
		fov:             fov,
		fovRays:         fovRays,
		viewingDistance: 30,
		brain:           brain,
		sprite:          sprite,
	}
}

func (a *Animal) IncrementPos(dx, dy float64) {
	a.x += dx
	a.y += dy
}

func (a Animal) GetPos() (x, y float64) {
	x = a.x
	y = a.y
	return x, y
}

func (a Animal) GetTheta() float64 {
	return a.dirTheta
}

func (a Animal) GetSpeed() float64 {
	return a.speed
}

func (a Animal) GetSprite() *pixel.Sprite {
	return a.sprite
}

func (a *Animal) turnDelta(dAngle float64) {
	a.dirTheta += dAngle
}

func (a *Animal) SetTurningState(x TurningState) {
	a.turningState = x
}

func (a *Animal) updateDirection() {
	a.see()
	a.brain.Think()
	turnLeft := a.brain.GetOutput(0) > 0.75
	turnRight := a.brain.GetOutput(1) > 0.75
	noTurn := a.brain.GetOutput(2) > 0.75
	if turnLeft {
		a.turningState = LEFT
	}
	if turnRight {
		a.turningState = RIGHT
	}
	if noTurn {
		a.turningState = STRAIGHT
	}

	switch a.turningState {
	case LEFT:
		a.turnDelta(a.turningRate)
	case RIGHT:
		a.turnDelta(-a.turningRate)
	}
}

func (a *Animal) Update() {
	a.updateDirection()
	dx := math.Cos(a.dirTheta) * a.speed * float64(1/float64(Game.TicksPerSecond))
	dy := math.Sin(a.dirTheta) * a.speed * float64(1/float64(Game.TicksPerSecond))
	a.x = a.x + dx
	a.y = a.y + dy

	a.ticksUntilHurt--
	if a.ticksUntilHurt <= 0 {
		a.hp--
		a.ticksUntilHurt = 20
	}

	if IsCellFull(a.x, a.y) {
		sx := a.x - math.Mod(a.x, 16.0)
		sy := a.y - math.Mod(a.y, 16.0)
		food := FoodBlocks[HashCoords(sx, sy)]
		if food.Eeat() > 0 {
			a.ticksUntilHurt += 100
			a.hp++
		}
	} else if IsCellFull(a.x+16, a.y) {
		sx := a.x+16 - math.Mod(a.x, 16.0)
		sy := a.y - math.Mod(a.y, 16.0)
		food := FoodBlocks[HashCoords(sx, sy)]
		if food.Eeat() > 0 {
			a.ticksUntilHurt += 100
			a.hp++
		}
	} else if IsCellFull(a.x, a.y+16) {
		sx := a.x - math.Mod(a.x, 16.0)
		sy := a.y+16 - math.Mod(a.y, 16.0)
		food := FoodBlocks[HashCoords(sx, sy)]
		if food.Eeat() > 0 {
			a.ticksUntilHurt += 100
			a.hp++
		}
	} else if IsCellFull(a.x+16, a.y+16) {
		sx := a.x+16 - math.Mod(a.x, 16.0)
		sy := a.y+16 - math.Mod(a.y, 16.0)
		food := FoodBlocks[HashCoords(sx, sy)]
		if food.Eeat() > 0 {
			a.ticksUntilHurt += 100
			a.hp++
		}
	}
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
	vec.X = a.x + a.w / 2
	vec.Y = a.y + a.h / 2

	dir := pixel.ZV
	dir.X = math.Cos(theta) * a.w / 2
	dir.Y = math.Sin(theta) * a.w / 2

	for i := 1; i <= a.viewingDistance; i++ {
		vec = vec.Add(dir)
		x := vec.X - math.Mod(vec.X, 16.0)
		y := vec.Y - math.Mod(vec.Y, 16.0)
		pos := pixel.ZV
		pos.X = x
		pos.Y = y
		Squares = append(Squares, pos)
		if FoodBlocks[HashCoords(x, y)] != nil {
			a.brain.SetVisionInput(idx, float64(i) / float64(a.viewingDistance+1))
			return
		}
	}
	a.brain.SetVisionInput(idx, 1000_000_000)
}

func (a *Animal) GetHP() int {
	return a.hp
}

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
