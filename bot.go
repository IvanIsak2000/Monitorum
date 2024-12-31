package main

import (
	"fmt"
	"log"
	"monitorum/handlers"
	"monitorum/utils"
	"strings"
	"time"
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
	config, err := utils.GetConfig()
	if err != nil{
		log.Fatal(err)
	}
	botInstance := MyBot{
		Token: config.Token, 
		ModerId: config.ModerId}
	botInstance.Init()
	botInstance.SendAlert("")
	

	update, _ := botInstance.Bot.UpdatesViaLongPolling(nil)
	defer botInstance.Bot.StopLongPolling()

	for update := range update{
		if update.Message != nil{
			text := update.Message.Text

			// Удаление сообщение если чат под защитой
			if utils.IsProtect(update.Message.Chat.ID) {
				go time.Sleep(1 * time.Second)
				
				// Удаляем сообщение
				botInstance.Bot.DeleteMessage(&telego.DeleteMessageParams{
					ChatID: update.Message.Chat.ChatID(),
					MessageID: update.Message.MessageID})
				// Если автобан - баним отправителя
				if utils.IsAutoban(update.Message.Chat.ID){
					handlers.BanForAutoban(botInstance.Bot, &update)
				}
			}

			// Удаление системного уведомления если это вступлени/выход из чата
			//
			if update.Message.NewChatMembers != nil {
				for _, user := range update.Message.NewChatMembers{
					fmt.Printf(
						"👉 Новый участник %v в чате %v %v\n", user.ID, update.Message.From.ID, update.Message.Chat.FirstName)
					botInstance.Bot.DeleteMessage(
						&telego.DeleteMessageParams{
							ChatID: update.Message.Chat.ChatID(),
							MessageID: update.Message.MessageID,
						})
				} 
			}

			// Участник вышел 
			//
			if update.Message.LeftChatMember != nil{
				fmt.Printf(
					"👈 Участник %v вышел из %v %v\n", update.Message.LeftChatMember.ID, update.Message.Chat.ID, update.Message.Chat.FirstName)
				botInstance.Bot.DeleteMessage(&telego.DeleteMessageParams{
					ChatID: update.Message.Chat.ChatID(),
					MessageID: update.Message.MessageID,
				})
			}

			// Обработка защиты чата
			if strings.HasPrefix(text, "!protect"){
				autoban := false

				arguments := strings.Split(text, " ")
				if len(arguments) > 1 && arguments[1] == "on"{
					autoban = true
				}

				if utils.UserIsAdmin(botInstance.Bot, &update){
					handlers.Protect(botInstance.Bot, &update, autoban)
				} else {
					handlers.YouAreNotAdmin(botInstance.Bot, &update)
				}
			}
			
			// Обработка прочих команд
			switch text{
				case "!start": 
					handlers.Start(botInstance.Bot, update.Message)
				case "/start": 
					handlers.Start(botInstance.Bot, update.Message)
				case "!status":
						handlers.Status(botInstance.Bot, &update)
				case "!ban":  //  добавить аргумент userId
					if utils.UserIsAdmin(botInstance.Bot, &update){
						handlers.Ban(botInstance.Bot, &update)
					} else {
						handlers.YouAreNotAdmin(botInstance.Bot, &update)
					}
				case "!report":
					handlers.Report(botInstance.Bot, &update)
				default:
					if strings.HasPrefix(text, "!"){
						handlers.Help(botInstance.Bot, update.Message)
					}
			}
		}
	}
}