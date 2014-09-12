package seq

import (
	"encoding/json"
	"github.com/lightning/go/types"
)

type Pitch types.Pitch
type Gain types.Gain

// Note contains the information to play a single note
// in a pattern
type note struct {
	sample string                  `json:"sample"`
	number types.Pitch             `json:"pitch"`
	velocity types.Gain            `json:"velocity"`
}

func (this *note) Sample() string {
	return this.sample
}

func (this *note) Number() types.Pitch {
	return this.number
}

func (this *note) Velocity() types.Gain {
	return this.velocity
}

// create a new Note instance
func NewNote(sample string, number types.Pitch, velocity types.Gain) types.Note {
	n := new(note)
	n.sample = sample
	n.number = number
	n.velocity = velocity
	return n
}

// parse a note from a json object
func ParseNote(ba []byte) (types.Note, error) {
	n := new(note)
	ed := json.Unmarshal(ba, n)
	if ed != nil {
		return nil, ed
	}
	return n, nil
}
