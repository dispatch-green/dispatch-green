package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

var Channels [11]string

func ResetChannels() {
	for i, _ := range Channels {
		Channels[i] = ""
	}
	Channels[2] = "Normal Patrol"
	Channels[0] = GenerateRadioMessage()
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
		} else {
			i, err := strconv.Atoi(channel)
			if err != nil {
				panic(err)

			}
			Channels[i] = r.FormValue("channel_text")
		}
	}
	fmt.Fprint(w, GenerateRadioMessage())
}

func GenerateRadioMessage() string {
	message := "/311 \n📻 Radio status \n"
	for i, channel := range Channels {
		if channel != "" && i != 0 {
			message += fmt.Sprintf("📢 Ch%d: %s\n", i, channel)
		}
	}
	message += "\nIf any other channels are in use or have been closed please reply to this 311"
	Channels[0] = message
	return message
}
