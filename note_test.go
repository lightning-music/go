package lightning

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestDecodeNote(t *testing.T) {
	note := new(Note)
	data := bytes.NewBufferString(`{"sample":"foo.wav","number":61,"velocity":87}`)
	err := json.Unmarshal(data.Bytes(), note)
	if err != nil {
		t.Fatal(err)
	}
}

func TestEncodeNote(t *testing.T) {
	note := Note{"bar.wav", 39, 105}
	buf, err := json.Marshal(note)
	if err != nil {
		t.Fatal(err)
	}
	data := bytes.NewBufferString(`{"sample":"bar.wav","number":39,"velocity":105}`)
	if 0 != bytes.Compare(data.Bytes(), buf) {
		t.Fatalf("wrong note")
	}
}
