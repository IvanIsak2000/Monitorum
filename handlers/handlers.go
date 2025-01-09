package handlers

import (
	"fmt"
	"log"
	"monitorum/utils"
	"github.com/mymmrac/telego"
)

func Start(bot *telego.Bot, message *telego.Message) {
	if message.Chat.Type == "private"{
		msg := `
🤚 Приветствую
Я бот для защиты чата.

⚡ Могу:
- удалять системные сообщения(вход/выход юзера)
- банить нарушителей порядка
- включать АнтиРейд режим
- репортить админу о нарушителе

⚠️ Для работы просто добавьте
меня в чат с правами администратора.
`		
		param := telego.SendMessageParams{
			ChatID: message.Chat.ChatID(),
			Text: msg,
			ReplyParameters: &telego.ReplyParameters{
				MessageID: message.MessageID,
			},
		}
		bot.SendMessage(&param)

	} else {
		msg := `
🤝 Спасибо, я добавлен в чат.

⚡ Могу:
- удалять системные сообщения(вход/выход юзера)
- банить нарушителей порядка
- включать АнтиРейд режим
- репортить админу о нарушителе

⚠️ Для работы я должен иметь
права администратора. Если
они выданы - просто проигнорируйте 
это сообщение.

💡 Проверить права - !status

⚙️ Посмотреть команды - !help
`
		param := telego.SendMessageParams{
			ChatID: message.Chat.ChatID(),
			Text: msg,
			ReplyParameters: &telego.ReplyParameters{
				MessageID: message.MessageID,
			},
		}
		_, err := bot.SendMessage(&param)
		if err != nil{
			log.Fatal(err)
		}
	}
}

func Help(bot *telego.Bot, message *telego.Message) {
	if message.Chat.Type == "private"{
		msg := `
⚠️ Команды, начинающиеся с "!" доступны
только в группах.

В чате команда должна начинаться с "!".
Например, !ban
`
		param := telego.SendMessageParams{
			ChatID: message.Chat.ChatID(),
			Text: msg,
			ReplyParameters: &telego.ReplyParameters{
				MessageID: message.MessageID,
			},
		}
		_, err := bot.SendMessage(&param)
		if err != nil{
			fmt.Println(err)}
	} else {
		msg := `
⚙️ Доступные:
<code>!start</code> или /start - Запуск бота.

<code>!status</code> - Просмотр статуса доступа.

<code>!ban</code> - Ликвидация юзера. Применять только в ответ на сообщение(reply)
Доступно только админам и если я имею права админа.

<code>!report</code> - Вызвать модератора. Поможет админу заметить нарушителя.

<code>!protect [BAN MODE]</code> - Включить антирейд режим. Удаляет все отправленные сообщения за следующие 60 секунд.
аргумент BAN MODE:
on - включает автобан
off - отключает автобан

<code>!help</code> - Вызвать это меню.
		`
				param := telego.SendMessageParams{
					ChatID: message.Chat.ChatID(),
					Text: msg,
					ParseMode: telego.ModeHTML,
					ReplyParameters: &telego.ReplyParameters{
						MessageID: message.MessageID,
					},}
				_, err := bot.SendMessage(&param)
				if err != nil{
					fmt.Println(err)}}
	
}


func Ban(bot *telego.Bot, update *telego.Update) {
	// Проверка, что чат это группу или супергруппа
	if update.Message.Chat.Type != "group" && update.Message.Chat.Type != "supergroup"{
		bot.SendMessage(&telego.SendMessageParams{
			ChatID: update.Message.Chat.ChatID(),
			Text: "❌ Команду /ban можно использовать только в группе или супергруппе",
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},})
		return
	}

	// Проверка на то, что у бота есть права
	me, _ := bot.GetMe()
	if !utils.BotIsAdmin(bot, me.ID, update.Message.Chat.ID){
		bot.SendMessage(&telego.SendMessageParams{
			ChatID: update.Message.Chat.ChatID(),
			Text: "❌ У меня нет прав. Дайте мне права администратора",
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
		})
		return

	}



	// Проверка, что это именно ответ на сообщение 
	target := update.Message.ReplyToMessage
	if target == nil{
		bot.SendMessage(&telego.SendMessageParams{
			ChatID: update.Message.Chat.ChatID(),
			Text: "❌ Для ликвидации необходимо ответить на сообщение с !ban",
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
		})
		return
	}

	// Проверка, что бан не для самого пользователя
	if target.From.ID == update.Message.From.ID{
		bot.SendMessage(&telego.SendMessageParams{
			ChatID: update.Message.Chat.ChatID(),
			Text: "❌ Вы не можете забанить себя",
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
		})
		return
	}
	
	// Получение бота
	if target.From.ID == me.ID{
		bot.SendMessage(
			&telego.SendMessageParams{
				ChatID: update.Message.Chat.ChatID(),
				Text: "❌ Я не могу забанить себя",
				ReplyParameters: &telego.ReplyParameters{
					MessageID: update.Message.MessageID,
				},
			})
		return 
	}

	// Сам бан
	err := bot.BanChatMember(
		&telego.BanChatMemberParams{
			ChatID: update.Message.Chat.ChatID(),
			UserID: target.From.ID,
			RevokeMessages: true})
	if err != nil {
		bot.SendMessage(
			&telego.SendMessageParams{
				ChatID: update.Message.Chat.ChatID(),
				ReplyParameters: &telego.ReplyParameters{
					MessageID: update.Message.MessageID,
				},
			Text: fmt.Sprintf("❌ Не удалось ликвидировать пользователя %v", target.From.ID),
		})
	} else {
		text := fmt.Sprintf("✅ Пользователь %v успешно ликвидирован", target.From.ID)
		bot.SendMessage(&telego.SendMessageParams{
			ChatID: update.Message.Chat.ChatID(),
			Text: text,
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
		})
	}
}


func BanForAutoban(bot *telego.Bot, update *telego.Update) {
	// Проверка, что чат это группу или супергруппа
	if update.Message.Chat.Type != "group" && update.Message.Chat.Type != "supergroup"{
		bot.SendMessage(&telego.SendMessageParams{
			ChatID: update.Message.Chat.ChatID(),
			Text: "❌ Команду /ban можно использовать только в группе или супергруппе",
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},})
		return
	}

	// Проверка на то, что у бота есть права
	me, _ := bot.GetMe()
	if !utils.BotIsAdmin(bot, me.ID, update.Message.Chat.ID){
		bot.SendMessage(&telego.SendMessageParams{
			ChatID: update.Message.Chat.ChatID(),
			Text: "❌ У меня нет прав. Дайте мне права администратора",
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
		})
		return

	}

	// Получение бота
	if update.Message.From.ID == me.ID{
		bot.SendMessage(
			&telego.SendMessageParams{
				ChatID: update.Message.Chat.ChatID(),
				Text: "❌ Я не могу забанить себя",
				ReplyParameters: &telego.ReplyParameters{
					MessageID: update.Message.MessageID,
				},
			})
		return 
	}

	// Сам бан
	err := bot.BanChatMember(
		&telego.BanChatMemberParams{
			ChatID: update.Message.Chat.ChatID(),
			UserID: update.Message.From.ID,
			RevokeMessages: true})
	if err != nil {
		bot.SendMessage(
			&telego.SendMessageParams{
				ChatID: update.Message.Chat.ChatID(),
				ReplyParameters: &telego.ReplyParameters{
					MessageID: update.Message.MessageID,
				},
			Text: fmt.Sprintf("❌ Не удалось ликвидировать пользователя %v", update.Message.From.ID),
		})
	} else {
		text := fmt.Sprintf("✅ Пользователь %v успешно ликвидирован", update.Message.From.ID)
		bot.SendMessage(&telego.SendMessageParams{
			ChatID: update.Message.Chat.ChatID(),
			Text: text,
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
		})
	}


}

func Protect(bot *telego.Bot, update *telego.Update, autoban bool) {
	if update.Message.Chat.Type != "group" && update.Message.Chat.Type != "supergroup"{
		bot.SendMessage(&telego.SendMessageParams{
			ChatID: update.Message.Chat.ChatID(),
			Text: "❌ Команду !protect можно использовать только в группе или супергруппе",
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
			
		})
		return
	}
	utils.AddChat(bot, update, autoban)

}

func YouAreNotAdmin(bot *telego.Bot, update *telego.Update) {
	bot.SendMessage(&telego.SendMessageParams{
		ChatID: update.Message.Chat.ChatID(),
		Text: "❌ Эту команду может использовать только админ",
		ReplyParameters: &telego.ReplyParameters{
			MessageID: update.Message.MessageID,
		},
	})
}

func Report(bot *telego.Bot, update *telego.Update) {
	if update.Message.Chat.Type != "group" && update.Message.Chat.Type != "supergroup"{
		bot.SendMessage(&telego.SendMessageParams{
			ChatID: update.Message.Chat.ChatID(),
			Text: "❌ Команду !report можно использовать только в группе или супергруппе",
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
		})
		return
	}

	// Проверка, что это именно ответ на сообщение 
	target := update.Message.ReplyToMessage
	me, _ := bot.GetMe()
	if target == nil{
		bot.SendMessage(&telego.SendMessageParams{
			ChatID: update.Message.Chat.ChatID(),
			Text: "❌ Для репорта необходимо ответить на сообщение с !report",
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
		})
		return
	}
	if target.From.ID == update.Message.From.ID{
		bot.SendMessage(&telego.SendMessageParams{
			ChatID: update.Message.Chat.ChatID(),
			Text: "❌ Вы не можете зарепортить себя",
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
		})
		return
	}
	
	// Получение бота
	if target.From.ID == me.ID{
		bot.SendMessage(
			&telego.SendMessageParams{
				ChatID: update.Message.Chat.ChatID(),
				Text: "❌ Я не могу зарепортить себя",
				ReplyParameters: &telego.ReplyParameters{
					MessageID: update.Message.MessageID,
				},})
		return 
	}

	text := fmt.Sprintf(`
🆘 Репорт из чата <b>%v</b>:
Пользователь <code>%v</code> пожаловался на пользователя <code>%v</code>
из-за: <blockquote>%v</blockquote>`, 
	update.Message.Chat.Title, 
	update.Message.From.ID,
	target.From.ID,
	update.Message.ReplyToMessage.Text,
)

	admins, _ := bot.GetChatAdministrators(&telego.GetChatAdministratorsParams{
		ChatID: update.Message.Chat.ChatID(),
	})
	for _, admin := range admins{
		_, err := bot.SendMessage(&telego.SendMessageParams{
			ChatID: telego.ChatID{ID: admin.MemberUser().ID},
			Text: text,
			ParseMode: telego.ModeHTML,
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
		})
		if err == nil {
			bot.SendMessage(&telego.SendMessageParams{
				ChatID: update.Message.Chat.ChatID(),
				Text: fmt.Sprintf("✅ Админ %v получил уведомление", admin.MemberUser().ID),
				ReplyParameters: &telego.ReplyParameters{
					MessageID: update.Message.MessageID,
				},
			})
			return 
		}
	}
	bot.SendMessage(&telego.SendMessageParams{
		ChatID: update.Message.Chat.ChatID(),
		Text: `
❌ Нет доступных админов. Админ должен быть админом 
в этой группе и должен запустить бота командой /start 
в личных сообщениях.`,
	ReplyParameters: &telego.ReplyParameters{
		MessageID: update.Message.MessageID,
	},
	})
}

func Status(bot *telego.Bot, update *telego.Update) {
	if update.Message.Chat.Type != "group" && update.Message.Chat.Type != "supergroup"{
		bot.SendMessage(&telego.SendMessageParams{
			ChatID: update.Message.Chat.ChatID(),
			Text: "❌ Команду !status можно использовать только в группе или супергруппе",
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
		})
		return
	}

	me, _ := bot.GetMe()
	if utils.BotIsAdmin(bot, me.ID, update.Message.Chat.ID){
		text := "✅ Права администратора выданы. Ничего дополнительного не требуется."
		bot.SendMessage(&telego.SendMessageParams{
			ChatID: update.Message.Chat.ChatID(),
			Text: text,
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
		})

	} else {
		text := "❌ Права администратора не выданы. Чат не защищён."

		bot.SendMessage(&telego.SendMessageParams{
			ChatID: update.Message.Chat.ChatID(),
			Text: text,
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
		})
	}
}