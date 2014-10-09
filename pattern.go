package lightning

import (
	"fmt"
)

// Tempo in bpm
type Tempo uint64

// Pattern encapsulates a sequence for a given sample
type Pattern struct {
	Length int      `json:"length"`
	Notes  [][]Note `json:"notes"`
}

// NotesAt returns a slice representing the notes
// that are stored at a particular position in a Pattern.
func (this *Pattern) NotesAt(pos Pos) []Note {
	notes := len(this.Notes)
	return this.Notes[int(pos)%notes]
}

// AddTo adds a Note to the Pattern at pos
func (this *Pattern) AddTo(pos Pos, note Note) error {
	var str string
	if int(pos) >= this.Length {
		str = "pos (%d) greater that pattern length (%d)"
		return fmt.Errorf(str, pos, this.Length)
	}
	if int(pos) < 0 {
		str = "pos (%d) less than 0"
		return fmt.Errorf(str, pos, this.Length)
	}
	this.Notes[int(pos)] = append(this.Notes[int(pos)], note)
	return nil
}

// NewPattern creates a Pattern with the specified
// initial size.
func NewPattern(size int) Pattern {
	return Pattern{
		size,
		make([][]Note, size),
	}
}
