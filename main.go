package main

import (
	"flag"
	"github.com/go-chi/chi"
	"net/http"
)

var (
	defaultPort = 3000
	defaultSettingsPath = "/etc/invoker/settings.json"

	settingsPath string
)

type settings struct {
	Port int `json:"port"`
	Commands map[string]string `json:"commands"`
}

func init() {
	flag.StringVar(&settingsPath, "settings", defaultSettingsPath, "")
	flag.Parse()
}

func main() {
	r := chi.NewRouter()
	r.Get("/invoke", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	http.ListenAndServe(":3000", r)
}





