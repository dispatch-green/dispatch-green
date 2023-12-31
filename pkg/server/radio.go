package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var Channels [11]string
var replymsg bool
var RadioHub *Hub

func CreateHub() {
	RadioHub = newHub()
}

func ResetChannels() {
	replymsg = true
	Channels = [11]string{GenerateRadioMessage(replymsg), "", "Normal Patrol", "", "", "", "", "", "", "", ""}
}

func LoadChannels() {
	channelFile, err := os.Open("./channels.json")
	if err != nil {
		return
	}

	byteValue, _ := ioutil.ReadAll(channelFile)

	err = json.Unmarshal(byteValue, &Channels)
	if err != nil {
		return
	}
}

func SaveChannels() {
	cJson, _ := json.Marshal(Channels)
	err := os.WriteFile("./channels.json", cJson, 0644)
	if err != nil {
		panic(err)
	}
}

func ProcessRadioUpdate(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			panic(err)
		}

		channel := strings.TrimPrefix(r.URL.Path, "/api/radio/")
		if channel == "reset" {
			ResetChannels()
			SaveChannels()
		} else if channel == "replymsg" {
			replymsg = !replymsg
		} else {
			i, err := strconv.Atoi(channel)
			if err != nil {
				panic(err)

			}
			Channels[i] = r.FormValue("channel_text")
			SaveChannels()
			message := "<textarea id=\"radioChannels\" name=\"radioChannels\" rows=\"10\" class=\"form-control\">\n"
			message += GenerateRadioMessage(replymsg)
			message += "\n</textarea>"
			RadioHub.broadcast <- []byte(message)
		}
	}
}

func GenerateRadioMessage(replymsg bool) string {
	message := "/311 \n📻 Radio status \n"
	for i, channel := range Channels {
		if channel != "" && i != 0 {
			message += fmt.Sprintf("📢 Ch%d: %s\n", i, channel)
		}
	}
	if replymsg {
		message += "\nIf any other channels are in use or have been closed please reply to this 311"
	}

	Channels[0] = message
	return message
}
