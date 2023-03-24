package main

import (
	"embed"
	"flag"
	"fmt"
	"github.com/mitchellh/go-homedir"
	logger2 "openai/logger"
	"openai/server"
	"os"
	"path/filepath"
)

var logger = logger2.GetLogger()

//go:embed dist
var webRoot embed.FS
var configFile string
var debugMode bool

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

	if configFile == "" {
		configFile = filepath.Join(configDir, "config.toml")
	}

	// start server
	s, err := server.NewServer(configFile)
	if err != nil {
		panic(err)
	}
	s.Run(webRoot, "dist", debugMode)
}

func init() {

	flag.StringVar(&configFile, "config", "", "Config file path (default: ~/.config/chat-gpt/config.toml)")
	flag.BoolVar(&debugMode, "debug", true, "Enable debug mode (default: true, recommend to set false in production env)")
	flag.Usage = usage
	flag.Parse()
}

func usage() {
	fmt.Printf(`WeChat-GPT, Version: 1.0.0
USAGE: 
  %s [command options]
OPTIONS:
`, os.Args[0])

	flagSet := flag.CommandLine
	order := []string{"config", "debug"}
	for _, name := range order {
		f := flagSet.Lookup(name)
		fmt.Printf("  --%s => %s\n", f.Name, f.Usage)
	}
}
