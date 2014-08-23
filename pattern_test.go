package lightning

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestNoteGain(t *testing.T) {
	note := NewNote(Gain(1.0), Pitch(1.0))
	assert.Equal(t, note.Gain(), Gain(1.0))
}

func TestNotePitch(t *testing.T) {
	note := NewNote(Gain(1.0), Pitch(0.5))
	assert.Equal(t, note.Pitch(), Pitch(0.5))
}
