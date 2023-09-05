package infrastructure

import (
	"context"
	"errors"
	"github.com/sashabaranov/go-openai"
	"sync"
)

// Models
//
// GPT432K0314             = "gpt-4-32k-0314"
// GPT432K                 = "gpt-4-32k"
// GPT40314                = "gpt-4-0314"
// GPT4                    = "gpt-4"
// GPT3Dot5Turbo0301       = "gpt-3.5-turbo-0301"
// GPT3Dot5Turbo           = "gpt-3.5-turbo"
// GPT3TextDavinci003      = "text-davinci-003"
// GPT3TextDavinci002      = "text-davinci-002"
// GPT3TextCurie001        = "text-curie-001"
// GPT3TextBabbage001      = "text-babbage-001"
// GPT3TextAda001          = "text-ada-001"
// GPT3TextDavinci001      = "text-davinci-001"
// GPT3DavinciInstructBeta = "davinci-instruct-beta"
// GPT3Davinci             = "davinci"
// GPT3CurieInstructBeta   = "curie-instruct-beta"
// GPT3Curie               = "curie"
// GPT3Ada                 = "ada"
// GPT3Babbage             = "babbage"

type ChatBot struct {
	apiKey string
	wg     *sync.WaitGroup
	client *openai.Client
	model  string
}

func NewChatBot(apiKey string, model string, wg *sync.WaitGroup) *ChatBot {
	client := openai.NewClient(apiKey)
	return &ChatBot{
		apiKey: apiKey,
		wg:     wg,
		client: client,
		model:  model,
	}
}

func (c *ChatBot) Ask(messages []openai.ChatCompletionMessage) (*string, error) {
	completionChan := make(chan openai.ChatCompletionResponse)
	errorsChan := make(chan error)
	c.wg.Add(1)
	go c.getCompletionResponse(messages, completionChan, errorsChan, c.wg)
	if <-errorsChan != nil {
		return nil, errors.New("was an error")
	}
	resp := <-completionChan
	content := resp.Choices[0].Message.Content

	return &content, nil
}

func (c *ChatBot) getCompletionResponse(messages []openai.ChatCompletionMessage, ch chan openai.ChatCompletionResponse, errCh chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    c.model,
			Messages: messages,
		},
	)

	if err != nil {
		errCh <- err
		return
	}
	ch <- resp
}
