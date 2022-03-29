package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type heist struct {
	status string
	name   string
}

var heists []heist
var HeistMessage string

func ResetHeists() {
	heists = []heist{
		{
			status: "❌",
			name:   "Bay City Bank",
		},
		{
			status: "❌",
			name:   "Fleeca GOHWY",
		},
		{
			status: "❌",
			name:   "Fleeca Harmony",
		},
		{
			status: "❌",
			name:   "Paleto Bank",
		},
		{
			status: "❌",
			name:   "Jewelry Store",
		},
		{
			status: "❌",
			name:   "Bobcat",
		},
		{
			status: "❌",
			name:   "Vault",
		},
		{
			status: "❌",
			name:   "Yacht",
		},
	}
}

func Process1090Update(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		err := r.ParseForm()
		if err != nil {
			panic(err)
		}

		heist := strings.TrimPrefix(r.URL.Path, "/api/1090/")
		i, err := strconv.Atoi(heist)
		if err != nil {
			panic(err)

		}

		for _, value := range r.PostForm {
			heists[i].status = value[0]
		}

	}
	fmt.Fprintf(w, Generate1090Message())
}

func Generate1090Message() string {
	HeistMessage := "/311 \n⏰ 10-90 o'Clock Tracker ⏰\n"
	for _, heist := range heists {
		HeistMessage += fmt.Sprintf("%s: %s\n", heist.status, heist.name)
	}
	return HeistMessage
}
