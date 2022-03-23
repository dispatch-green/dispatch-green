package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

var channels [11]string

func ProcessRadioUpdate(w http.ResponseWriter, r *http.Request) {
	channels[2] = "Main Comms"
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	channel := strings.TrimPrefix(r.URL.Path, "/api/radio/")
	i, err := strconv.Atoi(channel)
	if err != nil {
		panic(err)

	}
	channels[i] = r.FormValue("channel_text")

	message := "/311 📻 Radio status \n"
	for i, channel := range channels {
		if channel != "" {
			message += fmt.Sprintf("📡Ch%d: %s\n", i, channel)
		}
	}
	fmt.Fprintf(w, message)

}
