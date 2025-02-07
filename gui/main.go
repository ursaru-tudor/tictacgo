package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	WindowW int = 600
	WindowH int = 400
)

// Image assets
var (
	imageX    *ebiten.Image
	imageO    *ebiten.Image
	imageGrid *ebiten.Image
)

func loadImage(img **ebiten.Image, filename string) {
	var err error
	*img, _, err = ebitenutil.NewImageFromFile(filename)
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	loadImage(&imageX, "assets/X.png")
	loadImage(&imageO, "assets/O.png")
	loadImage(&imageGrid, "assets/grid.png")
}

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	screen.DrawImage(imageGrid, nil)
	screen.DrawImage(imageX, nil)

	var v ebiten.DrawImageOptions
	v.GeoM.Translate(64, 0)

	screen.DrawImage(imageO, &v)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WindowW / 2, WindowH / 2
}

func main() {
	ebiten.SetWindowSize(WindowW, WindowH)
	ebiten.SetWindowTitle("TicTacGo")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
