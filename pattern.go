package lightning

import (
	"errors"
)

// Note contains the information to play a single note
type Note interface {
	Gain() Gain
	Pitch() Pitch
}

// Pattern encapsulates a sequence for a given sample
type Pattern interface {
	Length() int
	NoteAt(pos Pos) Note
	AddNote(pos Pos, note Note)
}

// Note implementation
type noteImpl struct {
	gain Gain
	pitch Pitch
}

func (this *noteImpl) Gain() Gain {
	return this.gain
}

func (this *noteImpl) Pitch() Pitch {
	return this.pitch
}

// Pattern implementaion
type patternImpl struct {
	notes []Note
}

func (this *patternImpl) Length() int {
	return len(this.notes)
}

func (this *patternImpl) NoteAt(pos Pos) Note {
	notes := len(this.notes)
	if int(pos) > notes - 1 {
		panic(errors.New("step index out of bounds"))
	}
	return this.notes[pos]
}

func (this *patternImpl) AddNote(pos Pos, note Note) {
	this.notes = append(this.notes, note)
}

func NewNote(gain Gain, pitch Pitch) Note {
	return &noteImpl{
		gain,
		pitch,
	}
}

func NewPattern() Pattern {
	return &patternImpl{
		make([]Note, 16),
	}
}
