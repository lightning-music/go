package lightning

import (
	"fmt"
)

// tempo in bpm
type Tempo uint64

// Note contains the information to play a single note
type Note struct {
	Sample string          `json:"sample"`
	Number int             `json:"number"`
	Velocity int           `json:"velocity"`
}

// Pattern encapsulates a sequence for a given sample
type Pattern struct {
	Length int               `json:"length"`
	// notes array
	Notes [][]Note           `json:"notes"`
}

func (this *Pattern) NotesAt(pos Pos) []Note {
	notes := len(this.Notes)
	return this.Notes[ int(pos) % notes ]
}

func (this *Pattern) AddNote(pos Pos, note Note) error {
	var str string
	if int(pos) >= this.Length {
		str = "pos (%d) greater that pattern length (%d)"
		return fmt.Errorf(str, pos, this.Length)
	}
	if int(pos) < 0 {
		str = "pos (%d) less than 0"
		return fmt.Errorf(str, pos, this.Length)
	}
	this.Notes[ int(pos) ] = append(this.Notes[ int(pos) ], note)
	return nil
}

func NewNote(sample string, number int, velocity int) Note {
	return Note{ sample, number, velocity, }
}

// bardiv is a string of the form "1/<DIV>"
// where DIV can be any of
// 1, 2, 3, 4, 5, 6, 7, 8,
// 12, 16, 24, 32, 64, 128
func NewPattern(size int) Pattern {
	return Pattern{
		size,
		make([][]Note, size),
	}
}
