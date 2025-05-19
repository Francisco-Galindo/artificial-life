package game

import (
	"math"

	_ "image/png"

	"github.com/gopxl/pixel/v2"
)

type TurningState uint8

const (
	LEFT TurningState = iota
	RIGHT
	STRAIGHT
)

type Animal struct {
	x            float64
	y            float64
	nextX float64
	nextY float64
	dirTheta     float64
	speed        float64
	turningState TurningState
	turningRate  float64
	sprite *pixel.Sprite
}

func InitAnimal() *Animal {

	pic, err := LoadPicture("./img/looking-left.png")
	if err != nil {
		panic(err)
	}

	sprite := pixel.NewSprite(pic, pic.Bounds())

	return &Animal{speed: 30, turningRate: 0.05, sprite: sprite}
}

func (a *Animal) IncrementPos(dx, dy float64) {
	a.x += dx
	a.y += dy
}

func (a Animal) GetPos() (x, y float64){
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

func (a *Animal) Update() {
	switch a.turningState {
	case LEFT:
		a.turnDelta(a.turningRate)
	case RIGHT:
		a.turnDelta(-a.turningRate)
	}
	dx := math.Cos(a.dirTheta) * a.speed * float64(1 / Game.TicksPerSecond)
	dy := math.Sin(a.dirTheta) * a.speed * float64(1 / Game.TicksPerSecond)
	a.x = a.x + dx
	a.y = a.y + dy
	fmt.Printf("x: %v, y: %v\n", a.x, a.y)
}
