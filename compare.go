package rde

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func Diff(file1, file2 string) (string, error) {
	cmd := exec.Command("diff", "--speed-large-files", file1, file2)

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

// Compare finds diff between 2 txt files and extracts unique domains
func Compare(file1, file2 string) ([]string, error) {
	diff, err := Diff(file1, file2)
	if err != nil {
		return nil, err
	}

	diffSlice := strings.Split(diff, "\n")

	var domains []string
	domainsUnique := make(map[string]struct{})

	for _, line := range diffSlice {
		if line == "" || (!strings.HasPrefix(line, "<") && !strings.HasPrefix(line, ">")) {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) > 1 {
			domain := strings.TrimSuffix(parts[1], ".")
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
