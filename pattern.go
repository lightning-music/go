package lightning

import (
	"math"
)

// Note contains the information to play a single note
type Note struct {
	Gain Gain              `json:"gain"`
	Pitch Pitch            `json:"pitch"`
}

// Pattern encapsulates a sequence for a given sample
type Pattern struct {
	// the sample this pattern is associated with
	Sample string          `json:"sample"`
	// number of notes currently in the pattern
	Length int             `json:"length"`
	// notes array
	Notes []Note           `json:"notes"`
}

func (this *Pattern) NoteAt(pos Pos) Note {
	notes := len(this.Notes)
	return this.Notes[ int(pos) % notes ]
}

func (this *Pattern) AddNote(pos Pos, note Note) {
	if int(pos) >= len(this.Notes) {
		// expand Notes slice
		back := this.Notes
		this.Notes = make([]Note, nextPow2(int(pos)))
		copy(this.Notes, back)
	}
	this.Notes[ int(pos) ] = note
}

func (this *Pattern) AppendNote(note Note) {
	this.Notes = append(this.Notes, note)
}

func NewNote(gain Gain, pitch Pitch) Note {
	return Note{ gain, pitch, }
}

func NewPattern(sample string, sz int) Pattern {
	return Pattern{
		sample,
		sz,
		make([]Note, sz),
	}
}

// return the next highest power of 2
func nextPow2(i int) int {
	var exp int = 0
	for ; int(math.Pow(2, float64(exp))) <= i; exp++ {}
	return int(math.Pow(2, float64(exp)))
}
