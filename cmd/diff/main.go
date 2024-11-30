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

	comparator := rde.NewComparator(conf)
	comparator.Run()
}
