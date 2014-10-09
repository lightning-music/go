package main

import (
	"encoding/json"
	"flag"
	"github.com/hypebeast/go-osc/osc"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type Sample struct {
	Path string `json:"path,omitempty"`
}

// Trigger random OSC samples
func main() {
	httpPort := flag.Int("httpPort", 3428, "address of lightning http interface")
	oscPort := flag.Int("oscPort", 4800, "address of lightning osc interface")
	interval := flag.String("interval", "250ms", "time period between osc events")
	flag.Parse()
	// get the list of samples
	uri := "http://localhost:" + strconv.Itoa(*httpPort) + "/samples"
	resp, err := http.Get(uri)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var samples []Sample

	ba, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("could not read response body")
		log.Fatal(err)
	}

	err = json.Unmarshal(ba, &samples)
	if err != nil {
		log.Println("could not decode sample list")
		log.Fatal(err)
	}

	nsamples := len(samples)

	log.Printf("server reported %d samples\n", nsamples)

	oscClient := osc.NewOscClient("localhost", *oscPort)

	dur, err := time.ParseDuration(*interval)
	if err != nil {
		log.Println("could not parse duration")
		log.Fatal(err)
	}
	ticker := time.NewTicker(dur)
	for _ = range ticker.C {
		oscMsg := osc.NewOscMessage("/sample/play")
		oscMsg.Append(samples[ rand.Intn(nsamples) ].Path)
		oscMsg.Append(int32( rand.Intn(128) ))
		oscMsg.Append(int32( rand.Intn(128) ))
		err = oscClient.Send(oscMsg)
		if err != nil {
			log.Println("could not send osc message")
			log.Fatal(err)
		}
	}
}
