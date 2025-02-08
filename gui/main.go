package main

import (
	"bytes"
	_ "embed"
	"image"
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

// Image assets (embedded)

//go:embed assets/X.png
var imageDataX []byte

// Ebiten image data
var (
	imageX       *ebiten.Image
	imageO       *ebiten.Image
	imageGrid    *ebiten.Image
	imageTurn    *ebiten.Image
	imageWon     *ebiten.Image
	imageVictory *ebiten.Image
	imageDraw    *ebiten.Image
	imageRestart *ebiten.Image
)

const (
	GridEdgeLength = 64
)

func loadImageFile(img **ebiten.Image, filename string) {
	var err error
	*img, _, err = ebitenutil.NewImageFromFile(filename)
	if err != nil {
		log.Fatal(err)
	}
}

func loadImageEmbedded(img **ebiten.Image, imageData []byte) {
	imgDecoded, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		log.Fatal(err)
	}

	*img = ebiten.NewImageFromImage(imgDecoded)
}

func init() {
	loadImageFile(&imageX, "assets/X.png") //loadImageEmbedded(&imageX, imageDataX)
	loadImageFile(&imageO, "assets/O.png")
	loadImageFile(&imageGrid, "assets/grid.png")
	loadImageFile(&imageTurn, "assets/turn.png")
	loadImageFile(&imageWon, "assets/won.png")
	loadImageFile(&imageVictory, "assets/congrats.png")
	loadImageFile(&imageDraw, "assets/draw.png")
	loadImageFile(&imageRestart, "assets/restart.png")
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

func (g *Game) RestartGame() {
	g.current_turn = board.X
	g.ended = false
	g.gameBoard.Reset()
}

func ToleranceCheck(x int) bool {
	const ToleranceLimit = 3
	//x_adj := x / GridEdgeLength
	x_diff := x % GridEdgeLength

	if x_diff > GridEdgeLength/2 {
		x_diff = GridEdgeLength - x_diff
	}

	return x_diff <= ToleranceLimit
}

// Returns true if the point (x, y) is near a boundary between grid cells.
// Mitigates misinputs
func GridCellEdge(x, y int) bool {
	return ToleranceCheck(x) || ToleranceCheck(y)
}

func (g *Game) Update() error {

	// button to restart the game
	if g.ended {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			gridx := x / GridEdgeLength
			gridy := y / GridEdgeLength

			if gridy >= 2 && gridx >= 3 {
				g.RestartGame()
			}
		}
		return nil
	}

	if g.current_turn == board.NONE {
		g.current_turn = board.X
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if x < GridEdgeLength*3 && y < GridEdgeLength*3 {
			gridx := x / GridEdgeLength
			gridy := y / GridEdgeLength

			if GridCellEdge(x, y) {
				return nil
			}

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

	if g.ended {
		screen.DrawImage(imageRestart, GridPositionTranslated(board.Position{X: 3, Y: 2}))
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
