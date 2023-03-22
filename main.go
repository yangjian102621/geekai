package main

import (
	"embed"
	"github.com/mitchellh/go-homedir"
	logger2 "openai/logger"
	"openai/server"
	"os"
	"path/filepath"
)

var logger = logger2.GetLogger()

//go:embed dist
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

	if err != nil {
		logger.Errorf("failed to load web types: %v", err)
		return
	}

	// start server
	s, err := server.NewServer(filepath.Join(configDir, "/config.toml"))
	if err != nil {
		panic(err)
	}
	s.Run(webRoot, "dist")
}
