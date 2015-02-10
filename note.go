package lightning

// Note defines a MIDI-style note
type Note struct {
	Sample   string `json:"sample"`
	Number   int32  `json:"number"`
	Velocity int32  `json:"velocity"`
}
