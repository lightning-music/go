package lightning

import (
	"github.com/bmizerany/assert"
	"encoding/json"
	"testing"
)

func TestNoteGain(t *testing.T) {
	note := NewNote(Gain(1.0), Pitch(1.0))
	assert.Equal(t, note.Gain, Gain(1.0))
}

func TestNotePitch(t *testing.T) {
	note := NewNote(Gain(1.0), Pitch(0.5))
	assert.Equal(t, note.Pitch, Pitch(0.5))
}

func TestPatternLength(t *testing.T) {
	pat := NewPattern("foo.flac", 0)
	assert.Equal(t, pat.Length, 0)
}

func TestPatternNoteAt(t *testing.T) {
	pat := NewPattern("foo.flac", 0)
	pat.AddNote(0, NewNote(Gain(0.5), Pitch(2.0)))
	pat.AddNote(1, NewNote(Gain(1.0), Pitch(1.14)))
	pat.AddNote(2, NewNote(Gain(0.1), Pitch(1.26)))
	note := pat.NoteAt(2)
	assert.Equal(t, note.Pitch, Pitch(1.26))
}

func TestPatternAddNote(t *testing.T) {
	pat := NewPattern("foo.flac", 0)
	note := NewNote(Gain(0.5), Pitch(2.0))
	pat.AddNote(0, note)
}

func TestPatternAppendNote(t *testing.T) {
	pat := NewPattern("foo.flac", 0)
	pat.AppendNote(NewNote(Gain(0.01), Pitch(1.0)))
	pat.AppendNote(NewNote(Gain(0.02), Pitch(1.0)))
	pat.AppendNote(NewNote(Gain(0.03), Pitch(1.0)))
	pat.AppendNote(NewNote(Gain(0.04), Pitch(1.0)))
	assert.Equal(t, pat.NoteAt(3).Gain, Gain(0.04))
}

func TestNoteEncodeJson(t *testing.T) {
	expected := []byte(`{"gain":0.1,"pitch":2.5}`)
	bs, err := json.Marshal(NewNote(Gain(0.1), Pitch(2.5)))
	assert.Equal(t, err, nil)
	assert.Equal(t, bs, expected)
}

func TestNoteDecodeJson(t *testing.T) {
	actual := new(Note)
	expected := NewNote(Gain(0.2), Pitch(0.75))
	blob := []byte(`{"gain":0.2,"pitch":0.75}`)
	err := json.Unmarshal(blob, actual)
	assert.Equal(t, err, nil)
	assert.Equal(t, &expected, actual)
}

func TestPatternEncodeJson(t *testing.T) {
	pat := NewPattern("foo.flac", 1)
	pat.AddNote(0, NewNote(Gain(0.7), Pitch(0.2)))
	expected := []byte(`{"sample":"foo.flac","length":1,"notes":[{"gain":0.7,"pitch":0.2}]}`)
	bs, err := json.Marshal(pat)
	assert.Equal(t, err, nil)
	assert.Equal(t, bs, expected)
}

func TestPatternDecodeJson(t *testing.T) {
	expected := NewPattern("bar.flac", 2)
	expected.AddNote(0, NewNote(Gain(0.5), Pitch(1.2)))
	expected.AddNote(1, NewNote(Gain(0.7), Pitch(1.2)))
	bs := []byte(`{"sample":"bar.flac","length":2,"notes":[{"gain":0.5,"pitch":1.2},{"gain":0.7,"pitch":1.2}]}`)
	actual := new(Pattern)
	err := json.Unmarshal(bs, &actual)
	assert.Equal(t, err, nil)
	assert.Equal(t, &expected, actual)
}
