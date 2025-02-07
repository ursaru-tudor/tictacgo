package board

import (
	"errors"
	"fmt"
)

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

// Handles encountering bad players
type InvalidPlayerError struct {
	Value byte
}

func (e InvalidPlayerError) Error() string {
	return fmt.Sprintf("Invalid Player value, %v", e.Value)
}

func (p Player) String() string {
	if p == NONE {
		return "Nobody"
	} else if p == X {
		return "Player X"
	} else if p == O {
		return "Player O"
	} else { // Should this throw an error?
		return "Invalid Player"
	}
}

func AlternatePlayer(p Player) (Player, error) {
	switch p {
	case X:
		return O, nil
	case O:
		return X, nil
	case NONE:
		return p, errors.New("expected x or o, given none")
	default:
		return p, InvalidPlayerError{byte(p)}
	}
}

// Defines a position on the board
type Position struct {
	X, Y int
}

func (p Position) Valid() bool {
	return (p.X >= 0) && (p.X < 3) && (p.Y >= 0) && (p.Y < 3)
}

// Transforms indices from 1-indexed to 0-indexed
func (p *Position) Normalise() {
	p.X = p.X - 1
	p.Y = p.Y - 1
}

// Handles encountering bad positions
type OutOfBoundsPositionError struct {
	Where Position
}

func (e OutOfBoundsPositionError) Error() string {
	return fmt.Sprintf("Invalid Position object pointing to %v", e.Where)
}

// A board is defined by the owners of each position
type Board [3][3]Player

// Checks if a move is possible, returns an erorr if it can't determine
func (b Board) MovePossible(pos Position) (bool, error) {
	if !pos.Valid() {
		return false, OutOfBoundsPositionError{pos}
	}

	/* if !b[pos.X][pos.Y].Valid() {
	 *	 v := b[pos.X][pos.Y]
	 *	 return false, InvalidPlayerError{byte(v)}
	 * }
	 * Too much
	 */

	if b[pos.X][pos.Y] == NONE {
		return true, nil
	}

	return false, nil
}

// Actually changes the value at a position
func (b *Board) MarkPosition(pos Position, p Player) error {
	if !p.Valid() {
		return InvalidPlayerError{byte(p)}
	}

	v, err := b.MovePossible(pos)

	if err != nil {
		return err
	}

	if !v {
		return errors.New("position already occupied")
	}

	b[pos.X][pos.Y] = p

	return nil
}

func (b Board) OpenPositions() []Position {
	var positions []Position
	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			if b[x][y] == NONE {
				positions = append(positions, Position{X: x, Y: y})
			}
		}
	}
	return positions
}

func (b Board) CheckDraw() bool {
	return len(b.OpenPositions()) == 0 && b.GetWinner() == NONE
}

func (b Board) GetWinner() Player {
	// Check lines

	for line := 0; line < 3; line++ {
		if b[line][0] != NONE {
			good := true
			for i := 1; i < 3; i++ {
				if b[line][i] != b[line][0] {
					good = false
				}
			}
			if good {
				return b[line][0]
			}
		}
	}

	// Check columns
	for column := 0; column < 3; column++ {
		if b[0][column] != NONE {
			good := true
			for i := 1; i < 3; i++ {
				if b[i][column] != b[0][column] {
					good = false
				}
			}
			if good {
				return b[0][column]
			}
		}
	}

	// Check main diagonal

	good := true

	if b[0][0] != NONE {
		for i := 1; i < 3; i++ {
			if b[i][i] != b[0][0] {
				good = false
			}
		}
		if good {
			return b[0][0]
		}
	}

	// Check secondary diagonal

	good = true

	if b[0][2] != NONE {
		for i := 1; i < 3; i++ {
			if b[i][2-i] != b[0][2] {
				good = false
			}
		}

		if good {
			return b[0][2]
		}
	}

	return NONE
}
