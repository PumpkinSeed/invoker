package main

import (
	"flag"
	"github.com/PumpkinSeed/container-invoke/server"
)

var (
	defaultSettingsPath = "/etc/invoker/settings.json"

	settingsPath string
)

func init() {
	flag.StringVar(&settingsPath, "settings", defaultSettingsPath, "")
	flag.Parse()
}

func main() {
	server.Serve(settingsPath)
}





