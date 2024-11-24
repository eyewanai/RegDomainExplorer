package main

import (
	"log"
	"rde"
)

//func main() {
//	start := time.Now()
//	conf, err := rde.NewConf()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	files, err := rde.GetDirFiles(conf.Icaan.OutputFolder)
//	log.Printf("Found %d files", len(files))
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	previousFiles, err := rde.GetDirFiles(rde.GetPreviousDayPath(conf.Icaan.OutputFolder))
//	log.Printf("Found %d files", len(previousFiles))
//	if err != nil {
//		log.Println("Error: failed to perform diff.")
//		log.Println("You probably didn't download files from previous day.")
//		log.Println("If you have files from any date, you can just rename it to previous date.")
//		log.Fatal(err)
//	}
//
//	fileMap := make(map[string]string)
//	for _, file := range previousFiles {
//		_, filename := filepath.Split(file)
//		fileMap[filename] = file
//	}
//
//	output := "tmp.txt"
//	outputPrev := "tmp_prev.txt"
//
//	total := 0
//	totalNewDomains := 0
//
//	files = files[:40]
//
//	for _, file := range files {
//		_, currentFilename := filepath.Split(file)
//		previousFile, ok := fileMap[currentFilename]
//		if ok {
//			currentHash, err := rde.GetFileHash(file)
//			if err != nil {
//				log.Println("Error getting hash: ", err)
//				continue
//			}
//			previousHash, err := rde.GetFileHash(previousFile)
//			if err != nil {
//				log.Println("Error getting hash: ", err)
//				continue
//			}
//
//			if currentHash == previousHash {
//				log.Printf("Skipped diff for %s\n", currentFilename)
//				continue
//			} else {
//				if err := rde.Unzip(file, output); err != nil {
//					log.Fatalf("Error unzipping: ", err)
//				}
//				if err := rde.Unzip(previousFile, outputPrev); err != nil {
//					log.Fatalf("Error unzipping: ", err)
//				}
//
//				diff, err := rde.Compare(output, outputPrev)
//				if err != nil {
//					log.Fatalf("Error diff: ", err)
//				}
//
//				log.Printf("Found %d new domains for %s\n", len(diff), currentFilename)
//				totalNewDomains += len(diff)
//
//				if err := os.Remove(output); err != nil {
//					log.Fatalf("Error removing file: ", err)
//				}
//				if err := os.Remove(outputPrev); err != nil {
//					log.Fatalf("Error removing file: ", err)
//				}
//				total++
//			}
//			if total == 1500 {
//				os.Exit(1)
//			}
//		}
//	}
//
//	fmt.Printf("Total new domains: %d\n", totalNewDomains)
//	fmt.Printf("Finish processing for %d\n", time.Since(start))
//
//	//diffFilename := fmt.Sprintf("%s/daily_registered_%s.txt", conf.Icaan.BaseOutputFolder, time.Now().Format("2006-01-02"))
//	//
//	//rde.SaveToFile(diffFilename, []string{"google.com", "youtube.com"})
//
//}

func main() {
	conf, err := rde.NewConf()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	comparator := rde.NewComparator(conf)
	comparator.Run()
}
