package main

import (
	"embed"
	"github.com/mitchellh/go-homedir"
	logger2 "openai/logger"
	"openai/server"
	config2 "openai/types"
	"os"
	"path/filepath"
)

var logger = logger2.GetLogger()

//go:embed web
var webRoot embed.FS

func main() {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()

	// create config dir
	configDir, _ := homedir.Expand("~/.config/chat-gpt")
	_, err := os.Stat(configDir)
	if err != nil {
		err := os.MkdirAll(configDir, 0755)
		if err != nil {
			logger.Error(err)
			return
		}
	}

	// load service configs
	config, err := config2.LoadConfig(filepath.Join(configDir, "/config.toml"))
	if err != nil {
		logger.Errorf("failed to load web types: %v", err)
		return
	}

	// start server
	s := server.NewServer(config)
	s.Run(webRoot, "web")
}
