package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func ProcessReportUpdate(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			panic(err)
		}

		channel := strings.TrimPrefix(r.URL.Path, "/api/radio/")
		i, err := strconv.Atoi(channel)
		if err != nil {
			panic(err)

		}
		Channels[i] = r.FormValue("channel_text")

	}
	fmt.Fprintf(w, GenerateRadioMessage())
}

func GenerateReportMessage() string {
	message := "/311 📻 Radio status \n"
	for i, channel := range Channels {
		if channel != "" && i != 0 {
			message += fmt.Sprintf("📢Ch%d: %s\n", i, channel)
		}
	}
	Channels[0] = message
	return message
}
