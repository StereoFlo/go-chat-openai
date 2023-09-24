package entity

import (
	"github.com/sashabaranov/go-openai"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

type User struct {
	UserId     int64
	Messages   []openai.ChatCompletionMessage
	LastUpdate time.Time
}

type Users User

func (u *User) AddMessage(m openai.ChatCompletionMessage) {
	u.Messages = append(u.Messages, m)
}

func (u *User) IsTokenLimitReached() bool {
	fullStr := ""
	for _, v := range u.Messages {
		fullStr = fullStr + " " + v.Content
	}
	tokenizer := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}
	words := strings.FieldsFunc(fullStr, tokenizer)
	count := 0
	for _, word := range words {
		count += utf8.RuneCountInString(word)
	}
	return count >= 4096
}
