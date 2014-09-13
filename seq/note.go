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
	Samp string                  `json:"sample"`
	Num types.Pitch              `json:"number"`
	Vel types.Gain               `json:"velocity"`
}

func (this *note) Sample() string {
	return this.Samp
}

func (this *note) Number() types.Pitch {
	return this.Num
}

func (this *note) Velocity() types.Gain {
	return this.Vel
}

// create a new Note instance
func NewNote(sample string, number types.Pitch, velocity types.Gain) types.Note {
	n := new(note)
	n.Samp = sample
	n.Num = number
	n.Vel = velocity
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
