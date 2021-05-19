package cowin

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"os"
	"time"
)

var telegramBot *tgbotapi.BotAPI
var botCreateErr error

func init() {
	telegramBot, botCreateErr = tgbotapi.NewBotAPI("1897017778:AAECncT3H8vZ00aHCId6tVqetWK9P-sC8XY")
	if botCreateErr != nil {
		fmt.Println("error while creating telegram bot: ", botCreateErr)
		os.Exit(1)
	}
}

func MessageListener() {

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	for {
		calendarByDistrict, fetchErr := GetCalendarForDistrict("294")
		if fetchErr != nil {
			sendMessage(getChatIds(), "No vaccination centers are available")
		}
		covaxCenters := GetCovaxinCenters("covaxin", calendarByDistrict)
		availableCovaxCenters := GetAvailableCenters(covaxCenters)
		var replyContent string
		if len(availableCovaxCenters) == 0 {
			replyContent = "No vaccine centers available"
		} else {
			replyContent = fmt.Sprintf("Vaccines are available at %d centers. Details as follows\n", len(availableCovaxCenters))
			for i, availableCovaxCenter := range availableCovaxCenters {
				availableVaccines := 0
				for _, session := range availableCovaxCenter.Sessions {
					availableVaccines = availableVaccines + session.AvailableCapacity
				}
				replyContent = replyContent + fmt.Sprintf("%d) Center Name: %s, Available Slots: %d, PinCode: %d\n\n", i, availableCovaxCenter.Name, availableVaccines, availableCovaxCenter.Pincode)
			}
		}
		sendMessage(getChatIds(), replyContent)
		time.Sleep(15 * 60 * time.Second)
	}
}

func sendMessage(chatIds []int64, messageContent string) {
	for _, chatId := range chatIds {
		message := tgbotapi.NewMessage(chatId, messageContent)
		telegramBot.Send(message)
	}
}

func getChatIds() []int64 {
	return []int64{57380733, 54587111}
}
