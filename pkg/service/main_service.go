package service

import (
	"fmt"
	tgV5 "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go-chat-tg/pkg/infrastructure"
	"log"
	"sync"
)

type MainService struct {
	chatBot       *infrastructure.ChatBot
	chatBotApiKey string
	tgBot         *infrastructure.TelegramClient
	wg            *sync.WaitGroup
	userService   *UserService
}

func NewMainService(
	chatBot *infrastructure.ChatBot,
	tgBot *infrastructure.TelegramClient,
	wg *sync.WaitGroup,
	userMessages *UserService) *MainService {
	return &MainService{
		chatBot:     chatBot,
		tgBot:       tgBot,
		wg:          wg,
		userService: userMessages,
	}
}

func (ms *MainService) MessageHandler(update tgV5.Update) {
	_, user := ms.userService.GetUserById(update.Message.From.ID)
	if user.IsTokenLimitReached() {
		ms.ClearHistoryHandler(update)
	}
	ms.userService.AddUserMessage(update.Message.From.ID, &update.Message.Text)
	if user == nil {
		_, user = ms.userService.GetUserById(update.Message.From.ID)
	}
	res, err := ms.chatBot.Ask(user.Messages)
	if err != nil {
		log.Fatal(err)
	}
	ms.userService.AddAiMessage(user.UserId, res)
	ms.wg.Add(1)
	go ms.tgBot.SendMessage(update.Message.Chat.ID, res, user)
}

func (ms *MainService) HistoryHandler(update tgV5.Update) {
	history := ms.userService.GetHistory(update.Message.From.ID)
	if history == "" {
		history = "History is empty."
	}
	ms.wg.Add(1)
	_, user := ms.userService.GetUserById(update.Message.From.ID)
	go ms.tgBot.SendMessage(update.Message.Chat.ID, &history, user)
}

func (ms *MainService) ClearHistoryHandler(update tgV5.Update) {
	_, user := ms.userService.GetUserById(update.Message.From.ID)
	var m string
	if user.IsTokenLimitReached() {
		m = "Max message size limit of AI is reached. Your chat history has been cleared. Your previous message was not delivered:\n\n" + update.Message.Text
	} else {
		m = "Message history completely cleared."
	}
	ms.userService.ClearHistory(update.Message.From.ID)
	ms.wg.Add(1)
	go ms.tgBot.SendMessage(update.Message.Chat.ID, &m, user)
}

func (ms *MainService) StartHandler(update tgV5.Update) {
	_, user := ms.userService.GetUserById(update.Message.From.ID)
	m := fmt.Sprintf(ms.tgBot.GetWelcomeMessage(), update.SentFrom().LanguageCode)
	ms.userService.AddUserMessage(update.Message.From.ID, &m)
	res, err := ms.chatBot.Ask(user.Messages)
	if err != nil {
		log.Fatal(err)
	}
	user.Messages = nil
	ms.wg.Add(1)
	go ms.tgBot.SendMessage(update.Message.Chat.ID, res, nil)
}
