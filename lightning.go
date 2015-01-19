package lightning

// #cgo CFLAGS: -Wall -O2
// #cgo LDFLAGS: -L. -llightning -lm -ljack -lsndfile -lpthread -lsamplerate -logg
// #include <lightning/lightning.h>
// #include <lightning/types.h>
import "C"

import (
	"errors"
	"fmt"
	"math"
)

// lightning audio engine public interface
type Engine interface {
	// Connect JACK audio outputs
	Connect(ch1 string, ch2 string) error
	// Add a directory to the search list
	AddDir(file string) error
	// play an audio sample
	PlaySample(file string, pitch float64, gain float64) error
	// play a note
	PlayNote(note Note) error
	// start exporting
	ExportStart(file string) int
	// stop exporting
	ExportStop() int
}

// Note type
type Note interface {
	// Sample associated with the note
	Sample() string
	// Note number
	Number() int32
	// Velocity of the note
	Velocity() int32
}

type impl struct {
	handle C.Lightning
}

func (this *impl) Connect(ch1 string, ch2 string) error {
	err := int(C.Lightning_connect_to(this.handle, C.CString(ch1), C.CString(ch2)))
	if err != 0 {
		return errors.New("could not connect to JACK sinks")
	} else {
		return nil
	}
}

func (this *impl) AddDir(file string) error {
	if rc := int(C.Lightning_add_dir(this.handle, C.CString(file))); rc != 0 {
		return fmt.Errorf("Lightning_add_dir returned %d", rc)
	} else {
		return nil
	}
}

func (this *impl) PlaySample(file string, pitch float64, gain float64) error {
	err := C.Lightning_play_sample(
		this.handle, C.CString(file), C.pitch_t(pitch), C.gain_t(gain),
	)
	if err != 0 {
		return errors.New("could not play sample")
	} else {
		return nil
	}
}

func getPitch(note Note) float64 {
	return float64(math.Pow(2.0, (float64(note.Number())-60.0)/12.0))
}

func (this *impl) PlayNote(note Note) error {
	pitch := getPitch(note)
	gain := float64(float64(note.Velocity()) / 127.0)
	return this.PlaySample(note.Sample(), pitch, gain)
}

func (this *impl) ExportStart(file string) int {
	return int(C.Lightning_export_start(
		this.handle, C.CString(file),
	))
}

func (this *impl) ExportStop() int {
	return int(C.Lightning_export_stop(this.handle))
}

// Initialize a new lightning engine.
func NewEngine() Engine {
	instance := new(impl)
	instance.handle = C.Lightning_init()
	return instance
}
