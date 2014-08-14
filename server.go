
package main

import (
	"log"
	"net/http"
)

func main() {
	log.Fatal(http.ListenAndServe(":4200", http.FileServer(http.Dir("."))))
}

