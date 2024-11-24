package rde

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func runDiff(file1, file2 string) (string, error) {
	cmd := exec.Command("diff", "--speed-large-files", "-u", file1, file2)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return out.String(), nil
		}
		return "", fmt.Errorf("diff command failed: %s", stderr.String())
	}

	return out.String(), nil
}

func MapFilesByName(files []string) map[string]string {
	fileMap := make(map[string]string)
	for _, file := range files {
		_, filename := filepath.Split(file)
		fileMap[filename] = file
	}
	return fileMap
}

// HandleFileDiff handles the case where file hashes differ
func HandleFileDiff(curFile, prevFile, curFilename, tmpDir string) ([]string, error) {
	uniqueID := uuid.New().String()
	outputCur := fmt.Sprintf("%s/cur_%s.txt", tmpDir, uniqueID)
	outputPrev := fmt.Sprintf("%s/prev_%s.txt", tmpDir, uniqueID)

	if err := Unzip(curFile, outputCur); err != nil {
		return nil, fmt.Errorf("failed unzipping current file (%s): %v", curFile, err)
	}

	if err := Unzip(prevFile, outputPrev); err != nil {
		log.Printf("Error unzipping previous file (%s): %v\n", prevFile, err)
		return nil, fmt.Errorf("failed unzipping previous file (%s): %v", prevFile, err)
	}

	diff, err := Compare(outputCur, outputPrev)
	if err != nil {
		return nil, fmt.Errorf("failed comparing files (%s and %s): %v", outputCur, outputPrev, err)
	}

	log.Printf("Found %d new domains for %s\n", len(diff), curFilename)

	if err := os.Remove(outputCur); err != nil {
		log.Printf("Error removing file (%s): %v\n", outputCur, err)
	}
	if err := os.Remove(outputPrev); err != nil {
		log.Printf("Error removing file (%s): %v\n", outputPrev, err)
	}
	return diff, nil
}

// ProcessFileComparison - Probably we don't need compare hashes because 'diff' command do this for us
func ProcessFileComparison(curFile, prevFile, curFilename, tmpDir string) {
	curHash, err := GetFileHash(curFile)
	if err != nil {
		log.Printf("Error hashing current file (%s): %v\n", curFile, err)
		return
	}

	prevHash, err := GetFileHash(prevFile)
	if err != nil {
		log.Printf("Error hashing previous file (%s): %v\n", prevFile, err)
		return
	}

	if curHash != prevHash {
		HandleFileDiff(curFile, prevFile, curFilename, tmpDir)
	}
}

func Compare(file1, file2 string) ([]string, error) {
	diff, err := runDiff(file1, file2)
	if err != nil {
		return nil, err
	}

	diffSlice := strings.Split(diff, "\n")

	var domains []string
	domainsUnique := make(map[string]struct{})

	for _, line := range diffSlice {
		if line == "" || (!strings.HasPrefix(line, "-")) { //&& !strings.HasPrefix(line, ">")) {
			continue
		}

		line = line[1:]

		parts := strings.Fields(line)
		if len(parts) > 1 {
			domain := strings.TrimSuffix(parts[0], ".")
			if strings.HasPrefix(domain, "_") || strings.Count(domain, ".") != 1 {
				continue
			}
			if _, ok := domainsUnique[domain]; !ok {
				domainsUnique[domain] = struct{}{}
			}
		}
	}

	for domain := range domainsUnique {
		domains = append(domains, domain)
	}

	return domains, nil
}
