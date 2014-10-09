package lightning

// Sequencer provides a way to play a Pattern using timing
// events emitted from a Metro.
type Sequencer struct {
	metro   Metro
	pattern Pattern `json:"pattern"`
}

// NewSequencer creates a Sequencer.
func NewSequencer(patternSize int) *Sequencer {
	return new(Sequencer)
}

// AddTo adds a note to the Sequencer's pattern
// at pos.
func (this *Sequencer) AddTo(pos Pos, note Note) error {
	return this.pattern.AddTo(pos, note)
}

// Clear removes all the notes at a given position
// in the Sequencer's Pattern.
func (this *Sequencer) Clear(pos Pos) error {
	return this.pattern.Clear(pos)
}

// Start plays the Sequencer's Pattern.
func (this *Sequencer) Start() error {
	return this.metro.Start()
}

// Stop playing the Sequencer's Pattern.
func (this *Sequencer) Stop() {
	this.metro.Stop()
}
