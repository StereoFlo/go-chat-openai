package service

import (
	"errors"
	"fmt"
	tgV5 "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go-chat-tg/internal/entity"
	"go-chat-tg/pkg/openai"
	"go-chat-tg/pkg/telegram"
	"log"
	"sync"
)

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
	if user.IsTokenLimitReached() {
		ms.ClearHistoryHandler(update)
	}
	ms.userService.AddUserMessage(update.Message.From.ID, &update.Message.Text)
	if user == nil {
		_, user = ms.userService.GetUserById(update.Message.From.ID)
	}
	res, err := ms.chatBot.Ask(user.Messages)
	if err != nil {
		ms.sendMessageErrorHandler(update, err, user, nil)
		return
	}
	ms.userService.AddAiMessage(user.UserId, res)
	if 4096 <= len(*res) { //todo make constant instead
		messages := ms.splitString(*res, 4096)
		ms.wg.Add(1)
		go ms.tgBot.SendMany(update.Message.Chat.ID, messages, user)
		return
	}
	err = ms.tgBot.Send(update.Message.Chat.ID, res, user)
	if err != nil {
		ms.sendMessageErrorHandler(update, err, user, res)
	}
}

func (ms *MainService) HistoryHandler(update tgV5.Update) {
	history := ms.userService.GetHistory(update.Message.From.ID)
	if history == "" {
		history = "History is empty."
	}
	_, user := ms.userService.GetUserById(update.Message.From.ID)
	err := ms.tgBot.Send(update.Message.Chat.ID, &history, user)
	if err != nil {
		ms.sendMessageErrorHandler(update, err, user, nil)
	}
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
	err := ms.tgBot.Send(update.Message.Chat.ID, &m, user)
	if err != nil {
		ms.sendMessageErrorHandler(update, err, user, nil)
	}
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
	err = ms.tgBot.Send(update.Message.Chat.ID, res, nil)
	if err != nil {
		ms.sendMessageErrorHandler(update, err, user, nil)
	}
}

func (ms *MainService) splitString(str string, chunkSize int) []string {
	if len(str) == 0 {
		return nil
	}
	if chunkSize >= len(str) {
		return []string{str}
	}
	var chunks = make([]string, 0, (len(str)-1)/chunkSize+1)
	currentLen := 0
	currentStart := 0
	for i := range str {
		if currentLen == chunkSize {
			chunks = append(chunks, str[currentStart:i])
			currentLen = 0
			currentStart = i
		}
		currentLen++
	}
	chunks = append(chunks, str[currentStart:])
	return chunks
}

func (ms *MainService) sendMessageErrorHandler(update tgV5.Update, err error, user *entity.User, botRes *string) {
	log.Print(err)
	errMsg := "Sorry, service temporarily unavailable. Please try again later."
	err = ms.tgBot.Send(update.Message.Chat.ID, &errMsg, user)
	if err != nil {
		if botRes != nil {
			log.Print(errors.New(*botRes))
		}
		log.Print(err)
	}
}
