package entity

import (
	"github.com/sashabaranov/go-openai"
	"time"
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
