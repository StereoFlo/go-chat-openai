package service

import (
	"github.com/sashabaranov/go-openai"
	"go-chat-tg/internal/entity"
	"time"
)

type UserService struct {
	users []*entity.User
}

func NewUserService(users []*entity.User) *UserService {
	return &UserService{users: users}
}

func (us *UserService) GetUserMessages() []*entity.User {
	return us.users
}

func (us *UserService) ClearHistory(userId int64) {
	if nil == us.users {
		return
	}
	k, _ := us.GetUserById(userId)
	us.users = append(us.users[:k], us.users[k+1:]...)
}

func (us *UserService) GetUserById(userId int64) (int, *entity.User) {
	if len(us.users) <= 0 {
		user := us.registerUser(userId)
		return 0, user
	}
	for k, v := range us.users {
		if v.UserId == userId {
			return k, v
		}
	}
	user := us.registerUser(userId)
	return 0, user
}

func (us *UserService) AddUserMessage(userId int64, message *string) {
	_, user := us.GetUserById(userId)
	user.AddMessage(openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: *message,
	})
}

func (us *UserService) AddAiMessage(userId int64, message *string) {
	_, user := us.GetUserById(userId)
	user.AddMessage(openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: *message,
	})
}

func (us *UserService) GetHistory(userId int64) string {
	_, user := us.GetUserById(userId)
	if user == nil {
		return ""
	}
	str := ""
	for _, v := range user.Messages {
		str = str + v.Role + ": " + v.Content + "\n\n"
	}

	return str
}

func (us *UserService) registerUser(userId int64) *entity.User {
	user := &entity.User{
		UserId:     userId,
		Messages:   nil,
		LastUpdate: time.Time{},
	}
	us.users = append(us.users, user)
	return user
}
