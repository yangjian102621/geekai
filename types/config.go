package types

import (
	"bytes"
	"github.com/BurntSushi/toml"
	"net/http"
	"openai/utils"
	"os"
)

type Config struct {
	Listen   string
	Session  Session
	ProxyURL string
	OpenAi   OpenAi
}

// OpenAi configs struct
type OpenAi struct {
	ApiURL      string
	ApiKey      []string
	Model       string
	Temperature float32
	MaxTokens   int
}

// Session configs struct
type Session struct {
	SecretKey string // session encryption key
	Name      string
	Path      string
	Domain    string
	MaxAge    int
	Secure    bool
	HttpOnly  bool
	SameSite  http.SameSite
}

func NewDefaultConfig() *Config {
	return &Config{
		Listen: "0.0.0.0:5678",

		Session: Session{
			SecretKey: utils.RandString(64),
			Name:      "CHAT_GPT_SESSION_ID",
			Domain:    "",
			Path:      "/",
			MaxAge:    86400,
			Secure:    false,
			HttpOnly:  false,
			SameSite:  http.SameSiteLaxMode,
		},
		OpenAi: OpenAi{
			ApiURL:      "https://api.openai.com/v1/chat/completions",
			ApiKey:      []string{""},
			Model:       "gpt-3.5-turbo",
			MaxTokens:   1024,
			Temperature: 1.0,
		},
	}
}

func LoadConfig(configFile string) (*Config, error) {
	var config *Config
	_, err := os.Stat(configFile)
	if err != nil {
		config = NewDefaultConfig()
		// generate types file
		buf := new(bytes.Buffer)
		encoder := toml.NewEncoder(buf)

		if err := encoder.Encode(&config); err != nil {
			return nil, err
		}

		err := os.WriteFile(configFile, buf.Bytes(), 0644)
		if err != nil {
			return nil, err
		}

		return config, nil
	}
	_, err = toml.DecodeFile(configFile, &config)
	if err != nil {
		return nil, err
	}

	return config, err
}
