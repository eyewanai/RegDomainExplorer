package rde

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type Loader struct {
	Conf *Conf
}

func NewLoader(conf *Conf) *Loader {
	return &Loader{
		Conf: conf,
	}
}

func (l *Loader) Run() error {
	conf := l.Conf

	// Add unique folder name (date in this case)
	conf.Icaan.OutputFolder = fmt.Sprintf("%s/%s", conf.Icaan.OutputFolder, time.Now().Format("2006-01-02"))
	if err := os.MkdirAll(conf.Icaan.OutputFolder, 0755); err != nil {
		log.Fatal(err)
	}

	accessToken, err := Authenticate(conf)
	if err != nil {
		return err
	}

	links, err := GetZoneLinks(conf, accessToken)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	sem := make(chan struct{}, conf.Icaan.ConcurrentDownloads)

	for _, link := range links {
		sem <- struct{}{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := DownloadZone(link, accessToken, conf)
			if err != nil {
				log.Println(err)
			}
			<-sem
		}()
	}
	wg.Wait()

	return nil
}
