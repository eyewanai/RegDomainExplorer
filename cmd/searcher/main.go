package main

import (
	"log"
	"rde"
)

func main() {
	conf, err := rde.NewConf()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	searcher := rde.NewSeacher(conf)
	if err := searcher.Run(); err != nil {
		log.Fatal(err)
	}
}
