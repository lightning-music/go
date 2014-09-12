package seq

// Note contains the information to play a single note
// in a pattern
type note struct {
	Sample string          `json:"sample"`
	Number int             `json:"number"`
	Velocity int           `json:"velocity"`
}

// create a new Note instance
func NewNote(sample string, number int, velocity int) Note {
	return note{ sample, number, velocity, }
}
