// type definitions for lightning/go packages
package lightning

// this package should never import any other lightning packages
import (
	"net/http"
)

// binding to liblightning
type Engine interface {
	Connect(ch1 string, ch2 string) error
	AddDir(file string) int
	PlaySample(file string, pitch float64, gain float64) error
	PlayNote(note Note) error
	ExportStart(file string) int
	ExportStop() int
}

// a collection of samples that also manages its own
// endpoints in a rest api
type Samples interface {
	Samples() []string
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
