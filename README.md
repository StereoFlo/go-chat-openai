## Telegram bot

### Description

In this repository, there is a Telegram chatbot that uses OpenAI for communication. Its distinguishing feature is the history, which serves as context for the conversation. The history can also be cleared.

To run the bot, you need to:
1) create a bot in BotFather
2) get the token
3) register on openai.com
4) get the token
5) rename the .env.example file to .env and edit it

After that, you need to build the program, and it can be used.

### How to build

To build the bot, you need to install Go from the website https://go.dev/dl/, clone the repository, navigate to the repository folder, and execute the command `go build -o chat cmd/chat/main.go`. After that, you can run the binary file that was generated during the build from the command line.