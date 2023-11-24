package service

import (
	"fmt"
	tgV5 "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go-chat-tg/pkg/openai"
	"go-chat-tg/pkg/telegram"
	"sync"
)

const MaxMessageCharacters = 4096

type MainService struct {
	chatBot       *openai.ChatBot
	chatBotApiKey string
	tgBot         *telegram.TelegramClient
	wg            *sync.WaitGroup
	userService   *UserService
}

func NewMainService(
	chatBot *openai.ChatBot,
	tgBot *telegram.TelegramClient,
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
	ms.userService.AddUserMessage(update.Message.From.ID, &update.Message.Text)
	if user == nil {
		_, user = ms.userService.GetUserById(update.Message.From.ID)
	}
	res, err := ms.chatBot.Ask(user.Messages)
	if err != nil {
		errStr := err.Error()
		ms.wg.Add(1)
		go ms.tgBot.Send(update.Message.Chat.ID, &errStr, user)
		return
	}
	ms.userService.AddAiMessage(user.UserId, res)
	if MaxMessageCharacters <= len(*res) { //todo make constant instead
		messages := ms.splitString(*res, MaxMessageCharacters)
		ms.wg.Add(1)
		go ms.tgBot.SendMany(update.Message.Chat.ID, messages, user)
		return
	}
	ms.wg.Add(1)
	go ms.tgBot.Send(update.Message.Chat.ID, res, user)
}

func (ms *MainService) HistoryHandler(update tgV5.Update) {
	history := ms.userService.GetHistory(update.Message.From.ID)
	if history == "" {
		history = "History is empty."
	}
	ms.wg.Add(1)
	_, user := ms.userService.GetUserById(update.Message.From.ID)
	go ms.tgBot.Send(update.Message.Chat.ID, &history, user)
}

func (ms *MainService) ClearHistoryHandler(update tgV5.Update) {
	_, user := ms.userService.GetUserById(update.Message.From.ID)
	m := "Message history completely cleared."
	ms.userService.ClearHistory(update.Message.From.ID)
	ms.wg.Add(1)
	go ms.tgBot.Send(update.Message.Chat.ID, &m, user)
}

func (ms *MainService) StartHandler(update tgV5.Update) {
	_, user := ms.userService.GetUserById(update.Message.From.ID)
	m := fmt.Sprintf(ms.tgBot.GetWelcomeMessage(), update.SentFrom().LanguageCode)
	ms.userService.AddUserMessage(update.Message.From.ID, &m)
	res, err := ms.chatBot.Ask(user.Messages)
	if err != nil {
		errMsg := err.Error()
		ms.wg.Add(1)
		go ms.tgBot.Send(update.Message.Chat.ID, &errMsg, nil)
		return
	}
	user.Messages = nil
	ms.wg.Add(1)
	go ms.tgBot.Send(update.Message.Chat.ID, res, nil)
}

func (ms *MainService) splitString(s string, chunkSize int) []string {
	if len(s) == 0 {
		return nil
	}
	if chunkSize >= len(s) {
		return []string{s}
	}
	var chunks = make([]string, 0, (len(s)-1)/chunkSize+1)
	currentLen := 0
	currentStart := 0
	for i := range s {
		if currentLen == chunkSize {
			chunks = append(chunks, s[currentStart:i])
			currentLen = 0
			currentStart = i
		}
		currentLen++
	}
	chunks = append(chunks, s[currentStart:])
	return chunks
}
