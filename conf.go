package rde

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"time"
)

var (
	requiredEnv = []string{"ICAAN_USERNAME", "ICAAN_PASSWORD", "ICAAN_AUTH_URL", "ICAAN_BASE_URL", "ICAAN_OUTPUT_FOLDER"}
)

type Conf struct {
	Icaan IcaanConf
}

type IcaanConf struct {
	Username            string
	Password            string
	AuthURL             string
	BaseURL             string
	OutputFolder        string
	ConcurrentDownloads int
	DownloadRetries     int
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	fmt.Println(os.Getenv("icann_account_username"))
}

func getEnv(key string) (string, error) {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value, nil
	}
	return "", fmt.Errorf("env variable %s not set", key)
}

func NewConf() (*Conf, error) {
	envValues := make(map[string]string)

	for _, e := range requiredEnv {
		value, err := getEnv(e)
		if err != nil {
			return nil, err
		}
		envValues[e] = value
	}

	concurrentDownloads, err := strconv.Atoi(os.Getenv("ICAAN_CONCURRENT_DOWNLOADS"))
	if err != nil {
		return nil, err
	}

	downloadRetries, err := strconv.Atoi(os.Getenv("ICAAN_DOWNLOAD_RETRIES"))
	if err != nil {
		return nil, err
	}

	outputFolder := os.Getenv("ICAAN_OUTPUT_FOLDER")
	// Add unique folder name (date in this case)
	outputFolder = fmt.Sprintf("%s/%s", outputFolder, time.Now().Format("2006-01-02"))

	return &Conf{
		Icaan: IcaanConf{
			Username:            envValues["ICAAN_USERNAME"],
			Password:            envValues["ICAAN_PASSWORD"],
			AuthURL:             envValues["ICAAN_AUTH_URL"],
			BaseURL:             envValues["ICAAN_BASE_URL"],
			OutputFolder:        outputFolder,
			ConcurrentDownloads: concurrentDownloads,
			DownloadRetries:     downloadRetries,
		},
	}, nil
}
