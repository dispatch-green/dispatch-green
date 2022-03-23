package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"text/template"

	"github.com/ThePhilderbeast/NpQueueBot/pkg/server"
	"gopkg.in/yaml.v2"
)

var CONFIG_LOCATIONS = [2]string{
	"configs/config.yml",
	"config.yml"}

type Config struct {
	Port    int    `yaml:"port"`
	WebRoot string `yaml:"webroot"`
}

var cfg Config

func main() {

	cfg = loadConfig()

	// serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(cfg.WebRoot+"/web/static"))))

	http.HandleFunc("/api/radio/", server.ProcessRadioUpdate)

	// http.HandleFunc("/report", ServeReportPage)
	http.HandleFunc("/", ServePage)

	port := ":" + strconv.Itoa(cfg.Port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println("Failed to start server", err)
		return
	}

}

func loadConfig() Config {

	var cfg Config
	cfgParsed := false
	for _, loc := range CONFIG_LOCATIONS {
		f, err := os.Open(loc)
		if err == nil {
			decoder := yaml.NewDecoder(f)
			err = decoder.Decode(&cfg)
			if err != nil {
				panic(err)
			}
			cfgParsed = true
			fmt.Println("config read from " + loc)
			break
		}
		f.Close()
	}

	if !cfgParsed {
		panic("no config found")
	}
	return cfg

}

func parsePage(templatePath string) *template.Template {
	files := append([]string{templatePath}, getWidgetFiles()...)
	page, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	return page
}

func getWidgetFiles() []string {
	files, err := filepath.Glob(cfg.WebRoot + "/web/widgets/*.gohtml")
	if err != nil {
		panic(err)
	}
	return files
}

func ServePage(w http.ResponseWriter, r *http.Request) {
	tmpl := parsePage(cfg.WebRoot + "/web/pages/index.gohtml")
	tmpl.ExecuteTemplate(w, "index.gohtml", nil)
}
