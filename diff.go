package rde

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func runDiff(file1, file2 string) (string, error) {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("diff --speed-large-files -u %s %s | grep '^+' | sed 's/^+//' | grep -v '^++'", file1, file2))

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

func RunAwk(input, output string) error {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("awk '{print $1}' %s | sort -u | awk '{sub(/\\.$/, \"\"); print}' > %s", input, output))
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed run awk: %v", err)
	}
	err = os.Remove(input)
	if err != nil {
		return fmt.Errorf("failed delete file after awk: %v", err)
	}
	return nil
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

	outputCurClean := fmt.Sprintf("%s/cur_%s_clean.txt", tmpDir, uniqueID)
	outputPrevClean := fmt.Sprintf("%s/prev_%s_clean.txt", tmpDir, uniqueID)

	err := RunAwk(outputCur, outputCurClean)
	if err != nil {
		return nil, fmt.Errorf("failed run awk: %v", err)
	}

	err = RunAwk(outputPrev, outputPrevClean)
	if err != nil {
		return nil, fmt.Errorf("failed run awk: %v", err)
	}

	diff, err := runDiff(outputPrevClean, outputCurClean)
	if err != nil {
		return nil, fmt.Errorf("failed comparing files (%s and %s): %v", outputCur, outputPrev, err)
	}

	diffSlice := strings.Split(diff, "\n")

	log.Printf("Found %d new domains for %s\n", len(diffSlice), curFilename)

	if err := os.Remove(outputCurClean); err != nil {
		log.Printf("Error removing file (%s): %v\n", outputCur, err)
	}
	if err := os.Remove(outputPrevClean); err != nil {
		log.Printf("Error removing file (%s): %v\n", outputPrev, err)
	}
	return diffSlice, nil
}
