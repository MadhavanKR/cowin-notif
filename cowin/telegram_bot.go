package cowin

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
	"strings"
	"time"
)

var telegramBot *tgbotapi.BotAPI
var botCreateErr error
var chatIds []Subscription

type Subscription struct {
	name   string `json:"name"`
	chatId int64  `json:"chatId"`
}

func init() {
	telegramBot, botCreateErr = tgbotapi.NewBotAPI("1897017778:AAECncT3H8vZ00aHCId6tVqetWK9P-sC8XY")
	if botCreateErr != nil {
		log.Printf("error while creating telegram bot: %v\n", botCreateErr)
		os.Exit(1)
	}
	chatIds = make([]Subscription, 0)
}

func clearBackLogMesages(updates tgbotapi.UpdatesChannel) {
	updates.Clear()
}

func MessageListener() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	for {
		updates, fetchUpdateErr := telegramBot.GetUpdatesChan(u)
		if fetchUpdateErr != nil {
			log.Printf("error while fetching updates from telegram bot: %v\n", fetchUpdateErr)
			continue
		}
		clearBackLogMesages(updates)
		for update := range updates {
			if update.Message == nil {
				continue
			}

			if update.Message.IsCommand() {
				var replyContent string
				switch update.Message.Command() {
				case "help":
					replyContent = "type /subscribe to subscribe for vaccine notifications\n " +
						"/unsubscribe to stop receiving notifications\n" +
						"/list to see list of subscriptions"
					break
				case "subscribe":
					log.Printf("subscribing [%s] for vaccine updates", update.Message.Chat.FirstName)
					addChatId(update.Message.Chat.ID, update.Message.Chat.FirstName)
					replyContent = "Thanks for subscribing, you will start receiving notifications shortly"
					break
				case "unsubscribe":
					removeChatId(update.Message.Chat.ID)
					log.Printf("unsubscribing [%s] for vaccine updates", update.Message.Chat.FirstName)
					replyContent = "You've successfully unsubscribed from vaccine notifications"
					break
				case "list":
					usernames := getSubscribedUsernames()
					fmt.Println(usernames)
					replyContent = fmt.Sprintf("There are total of %d subscriptions, they are\n %v \n", len(usernames), usernames)
					break
				default:
					log.Printf("unknown command: %s\n", update.Message.Command())
					replyContent = "Unknown command. Please use /help to know about available commands"
					break
				}
				sendMessage([]Subscription{{
					name:   update.Message.Chat.UserName,
					chatId: update.Message.Chat.ID,
				}}, replyContent)
			}
		}
	}

}

func SendVaccineUpdates(vaccineCheckInterval int, vaccine string) {
	for {
		calendarByDistrict, fetchErr := GetCalendarForDistrict("294")
		if fetchErr != nil {
			sendMessage(getChatIds(), "Failed to query cowin apis")
			continue
		}
		covaxCenters := GetCovaxinCenters(vaccine, calendarByDistrict)
		availableCovaxCenters := GetAvailableCenters(covaxCenters)
		var replyContent string
		if len(availableCovaxCenters) == 0 {
			replyContent = "No vaccine centers available"
		} else {
			replyContent = fmt.Sprintf("Vaccines are available at %d centers. Details as follows\n", len(availableCovaxCenters))
			for i, availableCovaxCenter := range availableCovaxCenters {
				availableVaccines := 0
				sessionDetails := ""
				for _, session := range availableCovaxCenter.Sessions {
					if strings.ToLower(session.Vaccine) == vaccine {
						availableVaccines = availableVaccines + session.AvailableCapacity
						if session.AvailableCapacity > 0 {
							sessionDetails += fmt.Sprintf("date: %s slots: %v\n", session.Date, session.Slots)
						}
					}
				}
				if availableVaccines == 0 {
					replyContent = "No vaccine centers available"
				} else {
					replyContent = replyContent + fmt.Sprintf("%d) Center Name: %s, Available Slots: %d, PinCode: %d, Details: %s\n\n", i, availableCovaxCenter.Name, availableVaccines, availableCovaxCenter.Pincode, sessionDetails)
				}
			}
		}
		sendMessage(getChatIds(), replyContent)
		log.Printf("sent updates to %d chats\n", len(getChatIds()))
		time.Sleep(time.Duration(vaccineCheckInterval) * 60 * time.Second)
	}
}

func sendMessage(chatIds []Subscription, messageContent string) {
	for _, curSubscription := range chatIds {
		message := tgbotapi.NewMessage(curSubscription.chatId, messageContent)
		telegramBot.Send(message)
	}
}

func getSubscribedUsernames() []string {
	usernames := make([]string, 0)
	for _, curSubscription := range chatIds {
		usernames = append(usernames, curSubscription.name)
	}
	return usernames
}

func removeChatId(chatId int64) {
	index := -1
	for i, curSubscription := range chatIds {
		if chatId == curSubscription.chatId {
			index = i
			break
		}
	}
	if index != -1 {
		chatIds = append(chatIds[:index], chatIds[index+1:]...)
	}
}

func addChatId(chatId int64, username string) {
	for _, curSubscription := range chatIds {
		if curSubscription.name == username {
			return
		}
	}
	chatIds = append(chatIds, Subscription{
		name:   username,
		chatId: chatId,
	})
}

func getChatIds() []Subscription {
	return chatIds
}
