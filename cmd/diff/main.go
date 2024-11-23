package main

import (
	"log"
	"path/filepath"
	"rde"
)

func main() {
	conf, err := rde.NewConf()
	if err != nil {
		log.Fatal(err)
	}

	files, err := rde.GetDirFiles(conf.Icaan.OutputFolder)
	if err != nil {
		log.Fatal(err)
	}

	previousFiles, err := rde.GetDirFiles(conf.Icaan.OutputFolder)
	if err != nil {
		log.Fatal(err)
	}

	fileMap := make(map[string]string)
	for _, file := range previousFiles {
		_, filename := filepath.Split(file)
		fileMap[filename] = file
	}

	for _, file := range files {
		_, currentFilename := filepath.Split(file)
		previousFile, ok := fileMap[currentFilename]
		if ok {
			currentHash, err := rde.GetFileHash(file)
			if err != nil {
				log.Println("Error getting hash: ", err)
				continue
			}
			previousHash, err := rde.GetFileHash(previousFile)
			if err != nil {
				log.Println("Error getting hash: ", err)
				continue
			}

			if currentHash == previousHash {
				log.Printf("Skipped diff for %s\n", currentFilename)
				continue
			}
		}

	}

}
