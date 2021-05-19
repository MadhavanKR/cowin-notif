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
		covaxCenters := GetCovaxinCenters(calendarByDistrict)
		availableCovaxCenters := GetAvailableCenters(covaxCenters)
		var replyContent string
		if len(availableCovaxCenters) == 0 {
			replyContent = "No vaccine centers available"
		} else {
			centerNames := make([]int, 0)
			for _, availableCovaxCenter := range availableCovaxCenters {
				centerNames = append(centerNames, availableCovaxCenter.CenterId)
			}
			replyContent = fmt.Sprintf("Vaccines are available at %d centers. center ids %v", len(availableCovaxCenters), centerNames)
		}
		sendMessage(getChatIds(), replyContent)
		time.Sleep(60 * time.Second)
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
