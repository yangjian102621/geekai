package main

import (
	logger2 "chatplus/logger"
	"chatplus/server"
	"embed"
	"flag"
	"fmt"
	"os"
)

var logger = logger2.GetLogger()

//go:embed dist
var webRoot embed.FS
var configFile string
var debugMode bool

func main() {
	logger.Info("Loading config file: ", configFile)
	// start server
	s, err := server.NewServer(configFile)
	if err != nil {
		panic(err)
	}
	s.Run(webRoot, "dist", debugMode)
}

func init() {

	flag.StringVar(&configFile, "config", "config.toml", "Config file path (default: config.toml)")
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
