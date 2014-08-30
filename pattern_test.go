package lightning

import (
	"github.com/bmizerany/assert"
	"encoding/json"
	"testing"
)

func TestNoteVelocity(t *testing.T) {
	note := NewNote("audio/file.flac", 60, 120)
	assert.Equal(t, note.Velocity, 120)
}

func TestNoteNumber(t *testing.T) {
	note := NewNote("audio/file.flac", 60, 120)
	assert.Equal(t, note.Number, 60)
}

func TestNoteSample(t *testing.T) {
	note := NewNote("audio/file.flac", 60, 120)
	assert.Equal(t, note.Sample, "audio/file.flac")
}

func TestPatternLength(t *testing.T) {
	pat := NewPattern(0)
	assert.Equal(t, pat.Length, 0)
}

func TestPatternNoteAt(t *testing.T) {
	pat := NewPattern(0)
	pat.AddNote(0, NewNote("audio/file.flac", 60, 120))
	pat.AddNote(1, NewNote("audio/file.flac", 62, 120))
	pat.AddNote(2, NewNote("audio/file.flac", 64, 120))
	note := pat.NoteAt(2)
	assert.Equal(t, note.Number, 64)
}

func TestPatternAddNote(t *testing.T) {
	pat := NewPattern(0)
	note := NewNote("audio/file.flac", 72, 96)
	pat.AddNote(0, note)
}

func TestPatternAppendNote(t *testing.T) {
	pat := NewPattern(0)
	pat.AppendNote(NewNote("audio/file.flac", 48, 112))
	pat.AppendNote(NewNote("audio/file.flac", 50, 104))
	pat.AppendNote(NewNote("audio/file.flac", 52, 96))
	pat.AppendNote(NewNote("audio/file.flac", 54, 80))
	assert.Equal(t, pat.NoteAt(3).Velocity, 80)
}

func TestNoteEncodeJson(t *testing.T) {
	expected := []byte(`{"sample":"audio/file.flac","number":64,"velocity":108}`)
	bs, err := json.Marshal(NewNote("audio/file.flac", 64, 108))
	assert.Equal(t, err, nil)
	assert.Equal(t, bs, expected)
}

func TestNoteDecodeJson(t *testing.T) {
	actual := new(Note)
	expected := NewNote("audio/file.flac", 58, 109)
	blob := []byte(`{"sample":"audio/file.flac","number":58,"velocity":109}`)
	err := json.Unmarshal(blob, actual)
	assert.Equal(t, err, nil)
	assert.Equal(t, &expected, actual)
}

func TestPatternEncodeJson(t *testing.T) {
	pat := NewPattern(1)
	pat.AddNote(0, NewNote("audio/file.flac", 56, 101))
	expected := []byte(`{"length":1,"notes":[{"sample":"audio/file.flac","number":56,"velocity":101}]}`)
	bs, err := json.Marshal(pat)
	assert.Equal(t, err, nil)
	assert.Equal(t, bs, expected)
}

func TestPatternDecodeJson(t *testing.T) {
	expected := NewPattern(2)
	expected.AddNote(0, NewNote("audio/file.flac", 55, 84))
	expected.AddNote(1, NewNote("audio/file2.flac", 54, 76))
	bs := []byte(`{"length":2,"notes":[{"sample":"audio/file.flac","number":55,"velocity":84},{"sample":"audio/file2.flac","number":54,"velocity":76}]}`)
	actual := new(Pattern)
	err := json.Unmarshal(bs, &actual)
	assert.Equal(t, err, nil)
	assert.Equal(t, &expected, actual)
}
