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

//func main() {
//	filepath := "/Users/eyewan/Documents/Projects/DomainCrawler/icaan-domains/loader.txt.gz"
//	filepath2 := "/Users/eyewan/Documents/Projects/DomainCrawler/icaan-domains/loader.txt.gz"
//	out := "/Users/eyewan/Downloads/top.txt"
//	out2 := "/Users/eyewan/Downloads/org.txt"
//
//	start := time.Now()
//
//	err := rgx.Unzip(filepath, out)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	err = rgx.Unzip(filepath2, out2)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	log.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
//
//	start = time.Now()
//
//	_, err = rgx.Compare(out, out2)
//	if err != nil {
//		log.Fatal(err)
//	}
//	log.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
//}
