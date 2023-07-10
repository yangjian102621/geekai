package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

func HttpGet(uri string, proxy string) ([]byte, error) {
	var client *http.Client
	if proxy == "" {
		client = &http.Client{}
	} else {
		proxy, _ := url.Parse(proxy)
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxy),
			},
		}
	}

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func HttpPost(uri string, params map[string]interface{}, proxy string) ([]byte, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	var client *http.Client
	if proxy == "" {
		client = &http.Client{}
	} else {
		proxy, _ := url.Parse(proxy)
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxy),
			},
		}
	}

	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
