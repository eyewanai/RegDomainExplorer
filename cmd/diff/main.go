package main

import (
	"log"
	"rde"
	"time"
)

func main() {
	conf, err := rde.NewConf()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	st := time.Now()

	comparator := rde.NewComparator(conf)
	comparator.Run()

	log.Printf("Complete extracting daily registred domains for %d", time.Since(st))
}
