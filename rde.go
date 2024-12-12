package rde

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Loader struct {
	Conf         *Conf
	OutputFolder *string
}

func NewLoader(conf *Conf, outputFolder *string) *Loader {
	return &Loader{
		Conf:         conf,
		OutputFolder: outputFolder,
	}
}

func (l *Loader) Run() error {
	var outputFolder string

	conf := l.Conf

	if l.OutputFolder == nil {
		outputFolder = *l.OutputFolder
	} else {
		outputFolder = conf.Icaan.OutputFolder
	}

	if err := os.MkdirAll(outputFolder, 0755); err != nil {
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

type Comparator struct {
	Conf *Conf
}

func NewComparator(conf *Conf) *Comparator {
	return &Comparator{Conf: conf}
}

func (c *Comparator) Run() {
	conf := c.Conf

	// create .tmp dir for storing tmp txt files
	dirPath := ".tmp"

	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	curFiles, err := GetDirFiles(conf.Icaan.OutputFolder)
	if err != nil {
		log.Fatal(err)
	}

	prevFiles, err := GetDirFiles(GetPreviousDayPath(conf.Icaan.OutputFolder))
	if err != nil {
		log.Println("Error: failed to perform diff.")
		log.Println("You probably didn't download files from previous day.")
		log.Println("If you have files from any date, you can just rename it to previous date.")
		log.Fatal(err)
	}

	fileMap := MapFilesByName(prevFiles)

	// curFiles = curFiles[:30]

	wg := sync.WaitGroup{}
	diffTmpChannel := make(chan []string, len(curFiles))
	diffChannel := make(chan string, len(curFiles))

	for _, curFile := range curFiles {
		_, curFilename := filepath.Split(curFile)
		prevFile, ok := fileMap[curFilename]
		if ok {
			wg.Add(1)
			go func() {
				defer wg.Done()
				diff, err := HandleFileDiff(curFile, prevFile, curFilename, ".tmp")
				if err != nil {
					log.Printf("Error handling file diff for %s: %v\n", curFilename, err)
				} else {
					diffTmpChannel <- diff
				}

			}()
		}
	}

	go func() {
		wg.Wait()
		close(diffTmpChannel)
		for diff := range diffTmpChannel {
			for _, line := range diff {
				diffChannel <- line
			}
		}
		close(diffChannel)
	}()

	now := time.Now()
	filePath := fmt.Sprintf("%s/daily_registered_%d_%02d_%02d.txt", conf.Icaan.BaseOutputFolder, now.Year(), now.Month(), now.Day())

	err = WriteChannelToFile(diffChannel, filePath)
	if err != nil {
		log.Fatal(err)
	}
}
func WriteChannelToFile(ch <-chan string, filePath string) error {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for line := range ch {
		if line == "" {
			continue
		}
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("error writing to file: %v", err)
		}
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("error flushing to file: %v", err)
	}

	return nil
}

type Searcher struct {
	Conf *Conf
}

func NewSeacher(conf *Conf) *Searcher {
	return &Searcher{Conf: conf}
}

func (s *Searcher) Run() error {
	now := time.Now()
	filePath := fmt.Sprintf("%s/daily_registered_%d_%02d_%02d.txt", s.Conf.Icaan.BaseOutputFolder, now.Year(), now.Month(), now.Day())
	domains, err := ReadTxt(filePath)
	if err != nil {
		return err
	}

	reg, err := UnmarshalRegex()
	if err != nil {
		return err
	}

	Match(domains, reg)
	return nil
}
