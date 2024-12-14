package handlers

import (
	"fmt"
	"log"
	"monitorum/utils"
	"github.com/mymmrac/telego"
)

func Start(bot *telego.Bot, message telego.Message) {
	if message.Chat.Type == "private"{
		msg := `
🤚 Приветствую
Я бот для защиты чата.

Я могу банить нарушителей порядка,
репортить админу.

Для работы просто добавьте
меня в чат с правами администратора.
`		
		param := telego.SendMessageParams{
			ChatID: message.Chat.ChatID(),
			Text: msg,
		}
		bot.SendMessage(&param)

	} else {
		msg := `
Спасибо, я добавлен в чат.
Я могу банить нарушителей порядка,
репортить админу.

⚠️ Для работы я должен иметь
права администратора. Если
они выданы - просто проигнорируйте 
это сообщение.

Посмотреть команды - !help
`
		param := telego.SendMessageParams{
			ChatID: message.Chat.ChatID(),
			Text: msg,
		}
		_, err := bot.SendMessage(&param)
		if err != nil{
			log.Fatal(err)
		}
	}
}

func Help(bot *telego.Bot, message telego.Message) {
	if message.Chat.Type == "private"{
		msg := `
Команды, начинающиеся с "!" доступны
только в чате группы.

В чате команда должна начинаться с "!".
Например, !ban
`
		param := telego.SendMessageParams{
			ChatID: message.Chat.ChatID(),
			Text: msg,
		}
		_, err := bot.SendMessage(&param)
		if err != nil{
			fmt.Println(err)
		}
	} else {
		msg := `
Доступные:
!start или /start - Запуск бота

!ban - Ликвидация юзера. Применять только в ответ на сообщение(reply)
Доступно только админам и если я имею права админа.

!report - Вызвать модератора. Поможет админу заметить нарушителя.

!protect - Включить антирейд режим

!unprotect - Отключить антирейд режим

!help - Вызвать это меню.
		`
				param := telego.SendMessageParams{
					ChatID: message.Chat.ChatID(),
					Text: msg,
				}
				_, err := bot.SendMessage(&param)
				if err != nil{
					fmt.Println(err)
				}
	}
	
}


func Ban(bot *telego.Bot, update *telego.Update) {
	// Проверка, что чат это группу или супергруппа
	if update.Message.Chat.Type != "group" && update.Message.Chat.Type != "supergroup"{
		bot.SendMessage(&telego.SendMessageParams{
			ChatID: update.Message.Chat.ChatID(),
			Text: "❌ Команду /ban можно использовать только в группе или супергруппе",
		})
		return
	}

	me, _ := bot.GetMe()
	if !utils.CheckIsAdmin(bot, me.ID, update.Message.Chat.ID){
		bot.SendMessage(&telego.SendMessageParams{
			ChatID: update.Message.Chat.ChatID(),
			Text: "❌ У меня нет прав. Дайте мне права администратора",
		})
		return

	}


	// Проверка, что это именно ответ на сообщение 
	target := update.Message.ReplyToMessage
	if target == nil{
		bot.SendMessage(&telego.SendMessageParams{
			ChatID: update.Message.Chat.ChatID(),
			Text: "❌ Для ликвидации необходимо ответить на сообщение с /ban",
		})
		return
	}
	if target.From.ID == update.Message.From.ID{
		bot.SendMessage(&telego.SendMessageParams{
			ChatID: update.Message.Chat.ChatID(),
			Text: "❌ Вы не можете забанить себя",
		})
		return
	}
	
	// Получение бота
	if target.From.ID == me.ID{
		bot.SendMessage(
			&telego.SendMessageParams{
				ChatID: update.Message.Chat.ChatID(),
				Text: "❌ Я не могу забанить себя"},
		)
		return 
	}

	// Сам бан
	err := bot.BanChatMember(
		&telego.BanChatMemberParams{
			ChatID: update.Message.Chat.ChatID(),
			UserID: target.From.ID,
			RevokeMessages: true},
		)
	if err != nil {
		bot.SendMessage(
			&telego.SendMessageParams{
				ChatID: update.Message.Chat.ChatID(),
			Text: "❌ Не удалось ликвидировать пользователя. Возможно у вас недостаточно прав."},
		)
	} else {
		text := fmt.Sprintf("✅ Пользователь %v успешно ликвидирован", target.From.ID)
		bot.SendMessage(&telego.SendMessageParams{
			ChatID: update.Message.Chat.ChatID(),
			Text: text},
		)
	}
}
