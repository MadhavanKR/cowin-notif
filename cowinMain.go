package main

import (
	"flag"
	"fmt"
	"github.com/MadhavanKR/cowin-notif/cowin"
	"os"
	"strings"
)

func main() {
	var vaccineCheckInterval int
	var botToken string
	flag.IntVar(&vaccineCheckInterval, "interval", 1, "--interval=3")
	flag.StringVar(&botToken, "botToken", "default", "--botToken=<your-bot-token>")
	flag.Parse()
	if strings.Contains(botToken, "default") {
		fmt.Println("please supply botToken")
		os.Exit(1)
	}
	os.Setenv("botToken", botToken)
	go cowin.SendVaccineUpdates(vaccineCheckInterval, "covaxin")
	cowin.MessageListener()
}
