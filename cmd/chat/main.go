package main

import (
	"go-chat-tg/pkg/entity"
	"go-chat-tg/pkg/infrastructure"
	service2 "go-chat-tg/pkg/service"
	"os"
	"strings"
	"sync"
)

var historyBtnLabel = "History (%s)"
var clearHistoryBtnLabel = "Clear history"

func init() {
	//if err := godotenv.Load(); err != nil {
	//	log.Fatalln("no env gotten")
	//}
}

func main() {
	var wg sync.WaitGroup
	users := make([]*entity.User, 0)
	userService := service2.NewUserService(users)
	chatBot := infrastructure.NewChatBot(os.Getenv("OPENAI_API_KEY"), os.Getenv("AI_MODEL"), &wg)
	tgBot := infrastructure.NewTelegramClient(os.Getenv("TELEGRAM_API_KEY"), os.Getenv("WELCOME_MESSAGE"), historyBtnLabel, clearHistoryBtnLabel, &wg)
	service := service2.NewMainService(chatBot, tgBot, &wg, userService)
	updates := tgBot.GetUpdates()
	for update := range updates {
		if update.Message == nil {
			continue
		}
		switch true {
		case update.Message.Text == "/start":
			service.StartHandler(update)
			continue
		case strings.HasPrefix(update.Message.Text, "Clear history"):
			service.ClearHistoryHandler(update)
			continue
		case strings.HasPrefix(update.Message.Text, "History"):
			service.HistoryHandler(update)
			continue
		default:
			service.MessageHandler(update)
		}
		wg.Wait()
	}
}
