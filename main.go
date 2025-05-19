package main

import (
	"fmt"
	"image"
	"math"
	"os"
	"time"

	"example.com/artificial-life/game"

	_ "image/png"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"golang.org/x/image/colornames"
)

var animals []*game.Animal

func initProgram() {
	game.InitGameConfig()
	animals = append(animals, game.InitAnimal())
}

func loadPicture(path string) (pixel.Picture, error) {
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

func run() {
	cfg := opengl.WindowConfig{
		Title:  "Artificial Life",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := opengl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	defer win.Destroy()


	win.Clear(colornames.Skyblue)

	angle := 0.0
	last := time.Now()
	lastTick := time.Now()
	tickDuration := time.Duration(float64(1 / float64(game.Game.TicksPerSecond)) * float64(time.Second))

	for !win.Closed() {
		// dt := time.Since(last).Seconds()
		last = time.Now()

		win.Clear(colornames.Lightgreen)

		for _, animal := range(animals) {
			lol := float64(last.Sub(lastTick)) / float64(tickDuration)
			if last.Sub(lastTick) >= tickDuration {
				// fmt.Println("UPDATE", float64(1.0 / float64(game.Game.TicksPerSecond) * float64(time.Second)) )
				animal.Update()

				if win.Pressed(pixel.KeyLeft) {
					animal.SetTurningState(game.LEFT)
				} else if win.Pressed(pixel.KeyRight) {
					animal.SetTurningState(game.RIGHT)
				} else {
					animal.SetTurningState(game.STRAIGHT)
				}
				lastTick = last
			}
			fmt.Println(math.Round(lol*100))

			dx := math.Cos(animal.GetTheta()) * animal.GetSpeed() * float64(1 / float64(game.Game.TicksPerSecond)) * lol
			dy := math.Sin(animal.GetTheta()) * animal.GetSpeed() * float64(1 / float64(game.Game.TicksPerSecond)) * lol
			x, y := animal.GetPos()
			if lol > 0.95 {
			    fmt.Printf("x: %v, y: %v\n", x+dx, y+dy)
			}


			mat := pixel.IM
			mat = mat.ScaledXY(pixel.ZV, pixel.V(24, 24))
			mat = mat.Rotated(pixel.ZV, angle)
			mat = mat.Moved(pixel.Vec{X: x+dx, Y: y+dy})

			animal.GetSprite().Draw(win, mat)
		}


		win.Update()
	}
}

func main() {
	fmt.Println("Bienvenido")
	initProgram()
	opengl.Run(run)
}
