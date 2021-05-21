package main

import (
	"flag"
	"github.com/MadhavanKR/cowin-notif/cowin"
)

func main() {
	var vaccineCheckInterval int
	flag.IntVar(&vaccineCheckInterval,"interval", 1, "--interval=3")
	flag.Parse()
	go cowin.SendVaccineUpdates(vaccineCheckInterval, "covaxin")
	cowin.MessageListener()
}
