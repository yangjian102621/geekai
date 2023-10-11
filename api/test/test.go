package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	apiKey := "qjvqGdqpTY7qQaGBMenM7XgQ"
	apiSecret := "3G1RzBGXywZv4VbYRTyAfNns1vIOAG8t"
	token, err := getBaiduToken(apiKey, apiSecret)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(token)

}

func getBaiduToken(apiKey string, apiSecret string) (string, error) {

	url := fmt.Sprintf("https://aip.baidubce.com/oauth/2.0/token?client_id=%s&client_secret=%s&grant_type=client_credentials", apiKey, apiSecret)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error with send request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error with read response: %w", err)
	}
	var r map[string]interface{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return "", fmt.Errorf("error with parse response: %w", err)
	}

	if r["error"] != nil {
		return "", fmt.Errorf("error with api response: %s", r["error_description"])
	}

	return fmt.Sprintf("%s", r["access_token"]), nil
}
