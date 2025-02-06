package ttg

import "fmt"

type Player byte

const (
	NONE Player = 0
	X    Player = 'X'
	O    Player = 'O'
)

// Check if a Player variable is legal
// This shouldn't happen under normal circumstances
func (p Player) Valid() bool {
	return p == NONE || p == X || p == O
}

func (p Player) String() string {
	if !p.Valid() {
		return "Invalid Player"
	}

	if p == NONE {
		return "Nobody"
	}

	return fmt.Sprintf("%v", rune(p))
}
