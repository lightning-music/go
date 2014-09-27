package lightning

import (
	"encoding/json"
)

// Note contains the information to play a single note
// in a pattern
type Note struct {
	Samp string `json:"sample"`
	Num  int32  `json:"number"`
	Vel  int32  `json:"velocity"`
}

func (this *Note) Sample() string {
	return this.Samp
}

func (this *Note) Number() int32 {
	return this.Num
}

func (this *Note) Velocity() int32 {
	return this.Vel
}

// create a new Note instance
func NewNote(sample string, number int32, velocity int32) Note {
	return Note{sample, number, velocity}
}

// parse a note from a json object
func ParseNote(ba []byte) (*Note, error) {
	n := new(Note)
	ed := json.Unmarshal(ba, n)
	if ed != nil {
		return nil, ed
	}
	return n, nil
}
