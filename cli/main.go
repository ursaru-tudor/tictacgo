package main

import (
	"fmt"

	"github.com/ursaru-tudor/tictacgo/internal/board"
)

// Returns a string rendering a stylised arrangement of the game board
// Not optimised
func PresentBoard(b board.Board) string {
	var str string
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			switch b[i][j] {
			case board.NONE:
				str += " "
			case board.X:
				str += "X"
			case board.O:
				str += "O"
			}
			if j+1 < 3 {
				str += " | "
			}
		}
		str += "\n"

		for j := 0; i+1 < 3 && j < 3; j++ {
			str += "---"
		}

		if i+1 < 3 {
			str += "\n"
		}
	}
	return str
}

func main() {
	var game board.Board
	fmt.Println("Welcome to TicTacGo!")

	for current_player := board.X; game.GetWinner() == board.NONE && !game.CheckDraw(); current_player, _ = board.AlternatePlayer(current_player) {
		fmt.Printf("\n%v\n", PresentBoard(game))
		fmt.Printf("%v's turn: ", current_player)

		var pos board.Position

		fmt.Scan(&pos.X, &pos.Y)

		pos.Normalise()

		for can, err := game.MovePossible(pos); !can; can, err = game.MovePossible(pos) {
			if err != nil {
				fmt.Print("Invalid position, please try again\n")
			} else {
				fmt.Print("Position already occupied, please try again\n")
			}

			fmt.Scan(&pos.X, &pos.Y)

			pos.Normalise()
		}

		game.MarkPosition(pos, current_player)
	}

	fmt.Println(PresentBoard(game))

	if game.CheckDraw() {
		fmt.Printf("Nobody won...\n")
	} else {
		fmt.Printf("%v has won! Congratulations!!!\n", game.GetWinner())
	}
}
