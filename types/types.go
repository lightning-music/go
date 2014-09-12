// type definitions for lightning/go packages
package types

// this package should never import any other lightning packages
import (
	"github.com/gorilla/mux"
	"net/http"
)

// sugar
type Pitch float64
// sugar
type Gain float64

// thin wrapper around gorilla/mux/Router
type Router interface {
	Handle(path string, handler http.Handler) *Router
}

// Note
type Note interface {
	Sample() string
	Number() Pitch
	Velocity() Gain
}

// binding to liblightning
type Engine interface {
	AddDir(file string) int
	PlaySample(file string, pitch Pitch, gain Gain) error
	PlayNote(note Note) error
	ExportStart(file string) int
	ExportStop() int
}
