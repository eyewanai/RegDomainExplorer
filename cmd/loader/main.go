package main

import (
	"fmt"
	"log"
	"rde"
	"time"
)

func main() {
	start := time.Now()
	conf, err := rde.NewConf()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	loader := rde.NewLoader(conf)
	err = loader.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}
