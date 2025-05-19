package game

import (
	"image"
	"os"

	_ "image/png"

	"github.com/gopxl/pixel/v2"
)


type GameConfig struct {
	TicksPerSecond int
}

var Game GameConfig

func InitGameConfig() {
	Game = GameConfig{TicksPerSecond: 1}
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
