package server

import (
	"fmt"
	"net/http"
)

func ProcessReportUpdate(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			panic(err)
		}

		report := fmt.Sprintf("Title: %s | %s | %s \n", r.PostFormValue("reporterName"), r.PostFormValue("reportTitle"), r.PostFormValue("date"))
		report += "\n"
		report += "----\n"
		report += "\n"
		report += "Dispatch Report:\n"
		report += fmt.Sprintf("Report taken by: %s\n", r.PostFormValue("dispatcherName"))
		report += fmt.Sprintf("Date of Report: %s - %s\n", r.PostFormValue("date"), r.PostFormValue("time"))
		report += "\n"
		report += fmt.Sprintf("Date of Incident: %s - %s\n", r.PostFormValue("Incidentdate"), r.PostFormValue("Incidenttime"))
		report += "\n"
		report += "Reporting Person:\n"
		report += fmt.Sprintf("Name: %s\n", r.PostFormValue("reporterName"))
		report += fmt.Sprintf("SID: %s\n", r.PostFormValue("reporterSID"))
		report += fmt.Sprintf("Ph#: %s\n", r.PostFormValue("reporterPhNo"))
		report += "\n"
		if r.PostFormValue("Witness") != "" {
			report += "Person(s) Involved: :\n"
			report += r.PostFormValue("Witness") + "\n"
			report += "\n"
		}
		if r.PostFormValue("Suspect") != "" {
			report += "Suspect Description(s):\n"
			report += r.PostFormValue("Suspect") + "\n"
			report += "\n"
		}
		if r.PostFormValue("Vehicles") != "" {
			report += "Vehicle Description(s):\n"
			report += r.PostFormValue("Vehicles") + "\n"
			report += "\n"
		}
		if r.PostFormValue("location") != "" {
			report += "Location(s) of Incident:\n"
			report += r.PostFormValue("location") + "\n"
			report += "\n"
		}
		report += "Statement:\n"
		report += r.PostFormValue("Statement") + "\n"
		report += "\n"
		if r.PostFormValue("StolenItems") != "" {
			report += "Stolen Items:\n"
			report += r.PostFormValue("StolenItems") + "\n"
		}

		fmt.Fprint(w, report)
	}
}
