package main

import (
	"fmt"
	"log"
	"os"
	"rde"
	"sync"
	"time"
)

func main() {
	conf, err := rde.NewConf()
	if err != nil {
		log.Fatal("Error: ", err)
	}

	// Add unique folder name (date in this case)
	conf.Icaan.OutputFolder = fmt.Sprintf("%s/%s", conf.Icaan.OutputFolder, time.Now().Format("2006-01-02"))

	if err := os.MkdirAll(conf.Icaan.OutputFolder, 0755); err != nil {
		log.Fatal("Error: ", err)
	}

	accessToken, err := rde.Authenticate(conf)
	if err != nil {
		log.Fatal("Error: ", err)
	}

	links, err := rde.GetZoneLinks(conf, accessToken)

	// TODO: drop
	links = links[:10]

	if err != nil {
		log.Fatal("Error: ", err)
	}

	wg := sync.WaitGroup{}
	sem := make(chan struct{}, conf.Icaan.ConcurrentDownloads)

	for _, link := range links {
		sem <- struct{}{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := rde.DownloadZone(link, accessToken, conf)
			if err != nil {
				log.Println("Error: ", err)
			}
			<-sem
		}()
	}

	wg.Wait()

}

//func main() {
//	filepath := "/Users/eyewan/Documents/Projects/DomainCrawler/icaan-domains/dev.txt.gz"
//	filepath2 := "/Users/eyewan/Documents/Projects/DomainCrawler/icaan-domains/dev.txt.gz"
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
