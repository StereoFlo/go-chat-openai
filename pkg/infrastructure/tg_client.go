package infrastructure

import (
	"fmt"
	tgV5 "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go-chat-tg/pkg/entity"
	"log"
	"strconv"
	"sync"
)

type TelegramClient struct {
	bot                  *tgV5.BotAPI
	welcomeMessage       string
	wg                   *sync.WaitGroup
	historyBtnLabel      string
	clearHistoryBtnLabel string
}

func NewTelegramClient(
	apiKey string,
	welcomeMessage string,
	historyBtnLabel string,
	clearHistoryBtnLabel string,
	wg *sync.WaitGroup) *TelegramClient {
	bot, err := tgV5.NewBotAPI(apiKey)
	if err != nil {
		log.Panic(err)
	}
	return &TelegramClient{
		bot:                  bot,
		welcomeMessage:       welcomeMessage,
		wg:                   wg,
		historyBtnLabel:      historyBtnLabel,
		clearHistoryBtnLabel: clearHistoryBtnLabel,
	}
}
func (tc *TelegramClient) GetWelcomeMessage() string {
	return tc.welcomeMessage
}
func (tc *TelegramClient) GetUpdates() tgV5.UpdatesChannel {
	update := tgV5.NewUpdate(0)
	update.Timeout = 60
	return tc.bot.GetUpdatesChan(update)
}

func (tc *TelegramClient) SendMessage(userId int64, message *string, userMessage *entity.User) {
	defer tc.wg.Done()
	msg := tc.newMessage(userId, *message)
	var lens string
	if nil == userMessage || nil == userMessage.Messages {
		lens = "0"
	} else {
		lens = strconv.Itoa(len(userMessage.Messages))
	}
	tc.setKeyBoard(&msg, lens)
	_, err := tc.bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
}

func (tc *TelegramClient) setKeyBoard(msg *tgV5.MessageConfig, lens string) {
	msg.ReplyMarkup = tgV5.NewReplyKeyboard(
		tgV5.NewKeyboardButtonRow(
			tgV5.NewKeyboardButton(fmt.Sprintf(tc.historyBtnLabel, lens)),
			tgV5.NewKeyboardButton(tc.clearHistoryBtnLabel),
		),
	)
}

func (tc *TelegramClient) newMessage(userId int64, message string) tgV5.MessageConfig {
	return tgV5.MessageConfig{
		BaseChat: tgV5.BaseChat{
			ChatID:           userId,
			ReplyToMessageID: 0,
		},
		Text:                  message,
		DisableWebPagePreview: false,
		ParseMode:             "markdown",
	}
}
