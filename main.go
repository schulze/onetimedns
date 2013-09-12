package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"time"
)

var configFile = flag.String("config", "", "path to configuration file")

// Config represents a JSON config file ~/.config/onetimedns.
type Config struct {
	Name, Secret string
}

// Record describes the JSON data we get back from the OneTimeDNS server.
type Record struct {
	Expires int
	Address string `json:"record"`
	Status  string
	Value   string
}

func (r *Record) Get(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	dec := json.NewDecoder(resp.Body)
	if err = dec.Decode(r); err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func main() {
	flag.Parse()
	if len(*configFile) == 0 {
		home := os.Getenv("HOME")
		*configFile = home + "/.config/onetimedns"
	}

	file, err := os.Open(*configFile)
	if err != nil {
		log.Fatalln(err)
	}

	dec := json.NewDecoder(file)
	var c Config
	if err = dec.Decode(&c); err != nil {
		log.Fatalln(err)
	}
	url := "https://onetimedns.net/set?name=" + c.Name + "&secret=" + c.Secret

	// We should not to fail after this point.

	r := new(Record)
	tick := time.Tick(10 * time.Minute)
	for _ = range tick {
		err = r.Get(url)
		if err != nil {
			log.Println(err)
		} else {
			log.Println(r.Address)
		}
	}
}
