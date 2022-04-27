package utils

import "time"

func GetDate() string {

	est, _ := time.LoadLocation("SystemV/EST5EDT")
	date := time.Now().In(est).Format("02 Jan 2006")
	return date

}

func GetTime() string {

	est, _ := time.LoadLocation("SystemV/EST5EDT")
	time := time.Now().In(est).Format("15:04")
	return time

}
