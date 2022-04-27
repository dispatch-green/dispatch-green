package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/websocket"
	"gitlab.com/p6339/nopixel/dispatch-tools/pkg/server"
	"gitlab.com/p6339/nopixel/dispatch-tools/pkg/utils"
	"gopkg.in/yaml.v2"
)

var CONFIG_LOCATIONS = [2]string{
	"configs/config.yml",
	"config.yml"}

var upgrader = websocket.Upgrader{}

type Config struct {
	Port int `yaml:"port"`
}

type reportData struct {
	Date string
	Time string
}

var cfg Config

func main() {

	cfg = loadConfig()

	// load defaults
	server.ResetChannels()
	server.ResetHeists()
	server.CreateHub()

	// load channels if they are available
	server.LoadChannels()

	// serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(server.Static))))

	// api calls
	http.HandleFunc("/api/radio/", server.ProcessRadioUpdate)
	http.HandleFunc("/api/1090/", server.Process1090Update)
	http.HandleFunc("/api/report/", server.ProcessReportUpdate)
	http.HandleFunc("/api/radiows", wsUpgrade)

	// dynamic pages
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
	data := reportData{
		Date: utils.GetDate(),
		Time: utils.GetTime(),
	}
	tmpl.ExecuteTemplate(w, "index.gohtml", data)
}

func wsUpgrade(w http.ResponseWriter, r *http.Request) {
	log.Println("upgrading websocket")

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade failed: ", err)
		return
	}

	server.RegisterClient(ws, server.RadioHub)
}
