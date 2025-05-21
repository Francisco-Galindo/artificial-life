package main

import (
	"fmt"
	"image"
	"math"
	"math/rand"
	"os"
	"time"

	"example.com/artificial-life/game"

	_ "image/png"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"golang.org/x/image/colornames"
)

func initProgram() {
	game.InitGameConfig()
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

	var (
		camPos       = pixel.ZV
		camSpeed     = 500.0
		baseCamSpeed     = 500.0
		camZoom      = 1.0
		camZoomSpeed = 1.2
		// trees        []*pixel.Sprite
		// matrices     []pixel.Matrix
	)

	angle := 0.0
	last := time.Now()
	lastTick := time.Now()
	tickDuration := time.Duration(float64(1/float64(game.Game.TicksPerSecond)) * float64(time.Second))
	camPos.X = 2048
	camPos.Y = 2048

	for !win.Closed() {
		cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
		win.SetMatrix(cam)

		dt := time.Since(last).Seconds()
		last = time.Now()

		shouldUpdate := false
		if last.Sub(lastTick) >= tickDuration {
			lastTick = lastTick.Add(tickDuration)
			shouldUpdate = true
		}

		camZoom *= math.Pow(camZoomSpeed, win.MouseScroll().Y)
		camSpeed = baseCamSpeed / camZoom

		win.Clear(colornames.Forestgreen)

		if win.Pressed(pixel.KeyDown) {
			camPos.Y -= camSpeed * dt
		}
		if win.Pressed(pixel.KeyUp) {
			camPos.Y += camSpeed * dt
		}
		if win.Pressed(pixel.KeyLeft) {
			camPos.X -= camSpeed * dt
		}
		if win.Pressed(pixel.KeyRight) {
			camPos.X += camSpeed * dt
		}

		for _, food := range game.FoodBlocks {

			x, y := food.GetPos()
			mat := pixel.IM
			mat = mat.Rotated(pixel.ZV, angle)
			mat = mat.Moved(pixel.Vec{X: x, Y: y})

			food.GetCurrSprite().Draw(win, mat)
		}

		for _, animal := range game.Animals {

			var dx, dy float64
			if shouldUpdate {
				x := rand.Float64()
				if x < 0.33 {
					animal.SetTurningState(game.LEFT)
				} else if x < 0.75 {
					animal.SetTurningState(game.RIGHT)
				} else {
					animal.SetTurningState(game.STRAIGHT)
				}

				animal.Update()
			}

			tickProportion := float64(last.Sub(lastTick)) / float64(tickDuration)

			dx = math.Cos(animal.GetTheta()) * animal.GetSpeed() * float64(1/float64(game.Game.TicksPerSecond)) * tickProportion
			dy = math.Sin(animal.GetTheta()) * animal.GetSpeed() * float64(1/float64(game.Game.TicksPerSecond)) * tickProportion

			x, y := animal.GetPos()

			mat := pixel.IM
			mat = mat.Rotated(pixel.ZV, angle)
			mat = mat.Moved(pixel.Vec{X: x + dx, Y: y + dy})

			animal.GetSprite().Draw(win, mat)

			for _, papu := range game.Squares {
				mat := pixel.IM
				mat = mat.Rotated(pixel.ZV, angle)
				mat = mat.Moved(pixel.Vec{X: papu.X, Y: papu.Y})

				animal.GetSprite().Draw(win, mat)
			}
			game.Squares = make([]pixel.Vec, 0)
		}


		if shouldUpdate {
			game.PruneDeadAnimals()
		}

		win.Update()
	}
}

func main() {
	fmt.Println("Bienvenido")
	initProgram()
	opengl.Run(run)
}
