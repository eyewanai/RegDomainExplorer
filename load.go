package rde

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func doGet(url string, timeout int, headers map[string]string) (*http.Response, error) {

	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		r.Header.Add(k, v)
	}

	client := &http.Client{
		Timeout: time.Second * time.Duration(timeout),
	}
	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func GetZoneLinks(conf *Conf, accessToken string) ([]string, error) {
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", accessToken),
		"Content-Type":  "application/json",
		"Accept":        "application/json",
	}

	url := fmt.Sprintf("%s/czds/downloads/links", conf.Icaan.BaseURL)

	resp, err := doGet(url, 5, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var links []string
	if resp.StatusCode == http.StatusOK {
		if err := json.NewDecoder(resp.Body).Decode(&links); err != nil {
			return nil, err
		}
	}
	return links, nil
}

func DownloadZone(url string, accessToken string, conf *Conf) error {
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", accessToken),
		"Content-Type":  "application/json",
		"Accept":        "application/json",
	}

	start := time.Now()

	filename := fmt.Sprintf("%s.txt.gz", strings.Split(strings.Split(url, "/")[len(strings.Split(url, "/"))-1], ".")[0])
	fullPath := fmt.Sprintf("%s/%s", conf.Icaan.OutputFolder, filename)

	log.Printf("Downloading %s to %s\n", filename, conf.Icaan.OutputFolder)

	for attempts := 0; attempts < conf.Icaan.DownloadRetries; attempts++ {
		startByte := int64(0)
		if fileInfo, err := os.Stat(fullPath); err == nil {
			startByte = fileInfo.Size()
		}

		headers["Range"] = fmt.Sprintf("bytes=%d-", startByte)

		out, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Printf("Error opening file %s: %s\n", fullPath, err)
			return err
		}

		resp, err := doGet(url, 3600, headers)
		if err != nil {
			log.Printf("Error downloading %s: %s\n", url, err)
			out.Close()
			continue
		}

		defer resp.Body.Close()
		defer out.Close()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
			return fmt.Errorf("failed to download file, status code: %d\n", resp.StatusCode)
		}

		buf := make([]byte, 1024*1024*100) // 100 MB buffer, if you have stable connection you can try increase buf
		for {
			n, readErr := resp.Body.Read(buf)
			if n > 0 {
				if _, writeErr := out.Write(buf[:n]); writeErr != nil {
					log.Printf("Error writing to file %s: %s\n", fullPath, writeErr)
					return writeErr
				}
			}

			if readErr == io.EOF {
				log.Printf("Download completed for %s. Time spent %s\n", filename, time.Since(start))
				return nil
			}
			if readErr != nil {
				log.Printf("Error reading response body: %s\n", readErr)
				break
			}
		}
		log.Printf("Retrying download %s (%d/%d)\n", filename, attempts+1, conf.Icaan.DownloadRetries)
	}

	return fmt.Errorf("failed to download %s after %d attempts", filename, conf.Icaan.DownloadRetries)
}
