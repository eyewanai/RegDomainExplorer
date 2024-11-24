package rde

import (
	"compress/gzip"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

func GetDirFiles(dir string) ([]string, error) {
	var files []string

	items, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %v", err)
	}

	for _, item := range items {
		if !item.IsDir() {
			files = append(files, filepath.Join(dir, item.Name()))
		}
	}

	return files, nil
}

func GetFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file %s: %v", filePath, err)
	}
	defer file.Close()

	hash := md5.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return "", fmt.Errorf("failed to compute hash for file %s: %v", filePath, err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func GetPreviousDayPath(filePath string) string {
	_, datePart := filepath.Split(filePath)

	currentDate, err := time.Parse("2006-01-02", datePart)
	if err != nil {
		log.Println("Error parsing date:", err)
		return ""
	}

	previousDate := currentDate.AddDate(0, 0, -1)

	previousDateStr := previousDate.Format("2006-01-02")

	previousPath := filepath.Join(filepath.Dir(filePath), previousDateStr)
	return previousPath
}

func SaveToFile(filename string, lines []string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", filename, err)
	}
	defer file.Close()

	for _, line := range lines {
		_, err := file.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("failed to write to file %s: %v", filename, err)
		}
	}

	return nil
}

func Unzip(gzFilePath, outputFilePath string) error {
	gzFile, err := os.Open(gzFilePath)
	if err != nil {
		return err
	}
	defer gzFile.Close()

	gzReader, err := gzip.NewReader(gzFile)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	out, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, gzReader)
	return err
}

func CreateDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}
