package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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

func SymbolImage(p board.Player) (*ebiten.Image, error) {
	switch p {
	case board.X:
		return imageX, nil
	case board.O:
		return imageO, nil
	case board.NONE:
		return ebiten.NewImage(GridEdgeLength, GridEdgeLength), nil
	default:
		return nil, board.InvalidPlayerError{Value: byte(p)}
	}
}

type Game struct {
	gameBoard board.Board
}

func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if x <= GridEdgeLength*3 && y <= GridEdgeLength*3 {
			gridx := int(x / GridEdgeLength)
			gridy := int(y / GridEdgeLength)
			g.gameBoard[gridx][gridy] = board.X
		}
	}
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
			img, err := SymbolImage(g.gameBoard[i][j])

			// Should never be capable of happening
			if err != nil {
				panic(err)
			}

			screen.DrawImage(img, GridPositionTranslated(board.Position{X: i, Y: j}))
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
