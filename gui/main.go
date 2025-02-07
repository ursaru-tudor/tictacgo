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
	WindowW int = 320
	WindowH int = 196
)

// Image assets
var (
	imageX       *ebiten.Image
	imageO       *ebiten.Image
	imageGrid    *ebiten.Image
	imageTurn    *ebiten.Image
	imageWon     *ebiten.Image
	imageVictory *ebiten.Image
	imageDraw    *ebiten.Image
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
	loadImage(&imageTurn, "assets/turn.png")
	loadImage(&imageWon, "assets/won.png")
	loadImage(&imageVictory, "assets/congrats.png")
	loadImage(&imageDraw, "assets/draw.png")
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
		return nil, board.InvalidPlayerError(p)
	}
}

type Game struct {
	gameBoard    board.Board
	current_turn board.Player
	ended        bool
}

func (g *Game) Update() error {
	if g.ended {
		return nil
	}

	if g.current_turn == board.NONE {
		g.current_turn = board.X
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if x < GridEdgeLength*3 && y < GridEdgeLength*3 {
			gridx := int(x / GridEdgeLength)
			gridy := int(y / GridEdgeLength)

			moveposition := board.Position{X: gridx, Y: gridy}

			possible, _ := g.gameBoard.MovePossible(moveposition)

			if possible {
				g.gameBoard.MarkPosition(moveposition, g.current_turn)
				g.gameBoard[gridx][gridy] = g.current_turn
				g.current_turn, _ = board.AlternatePlayer(g.current_turn)
			}
		}
	}

	if g.gameBoard.CheckDraw() || g.gameBoard.GetWinner() != board.NONE {
		g.ended = true
		g.current_turn = board.NONE
	}

	return nil
}

func GridPositionTranslated(pos board.Position) *ebiten.DrawImageOptions {
	var dio ebiten.DrawImageOptions
	dio.GeoM.Translate(GridEdgeLength*float64(pos.X), GridEdgeLength*float64(pos.Y))
	return &dio
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Background
	screen.Fill(color.White)
	screen.DrawImage(imageGrid, nil)

	// Board
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			img, err := SymbolImage(g.gameBoard[i][j])

			// Should never be capable of happening
			if err != nil {
				log.Fatal(err)
			}

			screen.DrawImage(img, GridPositionTranslated(board.Position{X: i, Y: j}))
		}
	}

	// Game state
	if g.current_turn != board.NONE {
		img, _ := SymbolImage(g.current_turn)

		screen.DrawImage(img, GridPositionTranslated(board.Position{X: 3, Y: 0}))
		screen.DrawImage(imageTurn, GridPositionTranslated(board.Position{X: 4, Y: 0}))
	} else if g.gameBoard.GetWinner() != board.NONE {
		img, _ := SymbolImage(g.gameBoard.GetWinner())

		screen.DrawImage(img, GridPositionTranslated(board.Position{X: 3, Y: 0}))
		screen.DrawImage(imageWon, GridPositionTranslated(board.Position{X: 4, Y: 0}))
		screen.DrawImage(imageVictory, GridPositionTranslated(board.Position{X: 3, Y: 1}))
	} else {
		screen.DrawImage(imageDraw, GridPositionTranslated(board.Position{X: 3, Y: 1}))
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WindowW, WindowH
}

func main() {
	ebiten.SetWindowSize(WindowW*2, WindowH*2)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("TicTacGo")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
