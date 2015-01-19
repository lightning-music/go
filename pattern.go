package lightning

import (
	"fmt"
)

// Tempo in bpm
type Tempo uint64

// Pattern encapsulates a sequence for a given sample
type Pattern struct {
	Length     int      `json:"length"`
	Notes      [][]Note `json:"notes"`
}

type PatternEdit struct {
	Pos  Pos  `json:"pos"`
	Note Note `json:"note"`
}

func (this *Pattern) indexTooLarge(pos Pos) error {
	str := "pos (%d) greater than pattern length (%d)"
	return fmt.Errorf(str, pos, this.Length)
}

func (this *Pattern) indexNegative(pos Pos) error {
	return fmt.Errorf("pos (%d) less than 0", pos)
}

// NotesAt returns a slice representing the notes
// that are stored at a particular position in a Pattern.
func (this *Pattern) NotesAt(pos Pos) []Note {
	notes := len(this.Notes)
	return this.Notes[int(pos)%notes]
}

// AddTo adds a Note to the Pattern at pos
func (this *Pattern) AddTo(pos Pos, note Note) error {
	if int(pos) >= this.Length {
		return this.indexTooLarge(pos)
	}
	if int(pos) < 0 {
		return this.indexNegative(pos)
	}
	this.Notes[int(pos)] = append(this.Notes[int(pos)], note)
	return nil
}

// Clear removes all the notes at a given position
// in the Pattern.
func (this *Pattern) Clear(pos Pos) error {
	if int(pos) >= this.Length {
		return this.indexTooLarge(pos)
	}
	if int(pos) < 0 {
		return this.indexNegative(pos)
	}
	this.Notes[pos] = make([]Note, 0)
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
