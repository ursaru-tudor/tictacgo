package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/ursaru-tudor/tictacgo/internal/board"
)

// Screen sizes
const (
	WindowW int = 600
	WindowH int = 386
)

// Image assets
var (
	imageX    *ebiten.Image
	imageO    *ebiten.Image
	imageGrid *ebiten.Image
)

const (
	GridEdgeLength = 64
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

func GridPositionTranslated(pos board.Position) *ebiten.DrawImageOptions {
	var dio ebiten.DrawImageOptions
	dio.GeoM.Translate(GridEdgeLength*float64(pos.X), GridEdgeLength*float64(pos.Y))
	return &dio
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)

	screen.DrawImage(imageGrid, nil)

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			screen.DrawImage(imageO, GridPositionTranslated(board.Position{X: i, Y: j}))
		}
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WindowW / 2, WindowH / 2
}

func main() {
	ebiten.SetWindowSize(WindowW, WindowH)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("TicTacGo")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
