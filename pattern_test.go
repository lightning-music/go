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

func TestPatternLength(t *testing.T) {
	pat := NewPattern()
	assert.Equal(t, pat.Length(), 0)
}

func TestPatternAddNote(t *testing.T) {
	pat := NewPattern()
	note := NewNote(Gain(0.5), Pitch(2.0))
	pat.AddNote(0, note)
	assert.Equal(t, pat.Length(), 1)
}

func TestPatternNoteAt(t *testing.T) {
	pat := NewPattern()
	pat.AddNote(0, NewNote(Gain(0.5), Pitch(2.0)))
	pat.AddNote(1, NewNote(Gain(1.0), Pitch(1.14)))
	pat.AddNote(2, NewNote(Gain(0.1), Pitch(1.26)))
	note := pat.NoteAt(2)
	assert.Equal(t, note.Pitch(), Pitch(1.26))
}

func TestPatternNoteAppend(t *testing.T) {
	pat := NewPattern()
	pat.AppendNote(NewNote(Gain(0.01), Pitch(1.0)))
	pat.AppendNote(NewNote(Gain(0.02), Pitch(1.0)))
	pat.AppendNote(NewNote(Gain(0.03), Pitch(1.0)))
	pat.AppendNote(NewNote(Gain(0.04), Pitch(1.0)))
	assert.Equal(t, pat.Length(), 4)
	assert.Equal(t, pat.NoteAt(3).Gain(), Gain(0.04))
}
