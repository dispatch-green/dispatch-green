package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"

	"github.com/ThePhilderbeast/NpQueueBot/pkg/server"
	"gopkg.in/yaml.v2"
)

var CONFIG_LOCATIONS = [2]string{
	"configs/config.yml",
	"config.yml"}

type Config struct {
	Port int `yaml:"port"`
}

var cfg Config

func main() {

	cfg = loadConfig()

	server.ResetChannels()
	server.ResetHeists()

	// serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(server.Static))))

	http.HandleFunc("/api/radio/", server.ProcessRadioUpdate)
	http.HandleFunc("/api/1090/", server.Process1090Update)

	http.HandleFunc("/radio", ServeRadioPage)
	http.HandleFunc("/10-90", Serve1090Page)
	http.HandleFunc("/report", ServeReportPage)
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

func ServePage(w http.ResponseWriter, r *http.Request) {
	// tmpl, err := template.ParseFiles(cfg.WebRoot + "/web/pages/index.gohtml")
	tmpl, err := template.ParseFS(server.Templates, "web/pages/index.gohtml")
	if err != nil {
		panic(err)
	}
	tmpl.ExecuteTemplate(w, "index.gohtml", nil)
}

func ServeRadioPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFS(server.Templates, "web/pages/index.gohtml", "web/widgets/radio.gohtml")
	if err != nil {
		panic(err)
	}
	tmpl.ExecuteTemplate(w, "index.gohtml", server.Channels)
}

func Serve1090Page(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFS(server.Templates, "web/pages/index.gohtml", "web/widgets/10-90.gohtml")
	if err != nil {
		panic(err)
	}
	tmpl.ExecuteTemplate(w, "index.gohtml", server.Generate1090Message())
}

func ServeReportPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFS(server.Templates, "web/pages/index.gohtml", "web/widgets/report.gohtml")
	if err != nil {
		panic(err)
	}
	tmpl.ExecuteTemplate(w, "index.gohtml", server.Channels)
}
