package rde

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Authenticate return token
func Authenticate(conf *Conf) (string, error) {
	headers := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}
	credentials := map[string]string{
		"username": conf.Icaan.Username,
		"password": conf.Icaan.Password,
	}

	credentialsJson, _ := json.Marshal(credentials)

	r, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/api/authenticate", conf.Icaan.AuthURL),
		bytes.NewBuffer(credentialsJson),
	)
	if err != nil {
		log.Fatal(err)
	}

	for key, value := range headers {
		r.Header.Add(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var response map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			log.Fatal(err)
		}
		log.Println(response["message"])
		return response["accessToken"].(string), nil
	} else if resp.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("invalid authentication url %s", conf.Icaan.AuthURL)
	} else if resp.StatusCode != http.StatusUnauthorized {
		return "", fmt.Errorf("invalid username/password")
	} else if resp.StatusCode != http.StatusInternalServerError {
		return "", fmt.Errorf("internal server error. Please try again")
	} else {
		return "", fmt.Errorf("unknown error with code %d. Please try again", resp.StatusCode)
	}
}
