package main

import (
	"fmt"
	"monitorum/handlers"
	"os"
	"strings"
	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
)

type MyBot struct {
	Token string
	ModerId telego.ChatID
	Bot *telego.Bot
}

func (m *MyBot) Init() error {
	bot, err := telego.NewBot(m.Token)
	if err != nil{
		return fmt.Errorf("Ошибка создания бота: %w", err)
	}
	m.Bot = bot
	return nil
}

func (m *MyBot) SendAlert(text string) error {
	if text == ""{
		text = "✅ Бот запущен"
	}
	param := telego.SendMessageParams{ChatID: m.ModerId, Text: text}
	_, err := m.Bot.SendMessage(&param)
	if err != nil{
		return fmt.Errorf("Ошибка отправки сообщения: %w", err)
	}
	return nil
}

 
func main() {
    err := godotenv.Load()
    if err != nil {
        fmt.Println("Ошибка загрузки .env файла:", err)
    }
	t := os.Getenv("TOKEN")
	var m_id int64 =  5818860026
	botInstance := MyBot{
		Token: t, 
		ModerId: telego.ChatID{ID: m_id},
	}
	botInstance.Init()
	botInstance.SendAlert("")

	update, err := botInstance.Bot.UpdatesViaLongPolling(nil)
	if err != nil{fmt.Println(err)}
	defer botInstance.Bot.StopLongPolling()

	for update := range update{
		if update.Message != nil{
			text := update.Message.Text
			switch text{
				case "!start": 
					handlers.Start(botInstance.Bot, *update.Message)
				case "/start": 
					handlers.Start(botInstance.Bot, *update.Message)
				case "!ban": 
					handlers.Ban(botInstance.Bot, &update)
				default:
					if strings.HasPrefix(text, "!"){
						handlers.Help(botInstance.Bot, *update.Message)
					}
			}
		}
	}
}