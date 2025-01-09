package utils

import (
	"fmt"
	"github.com/mymmrac/telego"
	"time"
)

type ChatItem struct{
	IsProtected bool
	Autoban bool
}

var ProtectionChats map[int64]ChatItem

func init() {
	// Инициализация глобальной переменной ProtectionChats
	ProtectionChats = make(map[int64]ChatItem)
}


func AddChat(bot *telego.Bot, update *telego.Update, autoban bool) {
	ProtectionChats[update.Message.Chat.ID] = ChatItem{
		IsProtected: true,
		Autoban: autoban,
	}

	base_text := "🚨 ВКЛЮЧЕН АНТИРЕЙД РЕЖИМ 🚨\nВСЕ СООБЩЕНИЯ В ТЕЧЕНИИ 60 СЕКУНД БУДУТ УДАЛЯТЬСЯ"
	if autoban{
		base_text += " А ОТПРАВИТЕЛИ БАНИТЬСЯ"
	}

	notif_msg, _ := bot.SendMessage(
		&telego.SendMessageParams{
			ChatID: update.Message.Chat.ChatID(),
			Text: base_text,
			ReplyParameters: &telego.ReplyParameters{
				MessageID: update.Message.MessageID,
			},
		})
	fmt.Printf("Чат %v %v добавлен в защиту\n", update.Message.Chat.ID, update.Message.Chat.Title)
	
	go func() {
		time.Sleep(60 * time.Second)
		DeleteChat(bot, update, notif_msg)
	}()
}

func DeleteChat(bot *telego.Bot, update *telego.Update, notif_msg *telego.Message) {
	if _, exists := ProtectionChats[update.Message.Chat.ID]; exists{
		delete(ProtectionChats, update.Message.Chat.ID)

		// Удаление сообщение о предупреждении
		bot.DeleteMessage(&telego.DeleteMessageParams{
			ChatID: update.Message.Chat.ChatID(),
			MessageID: notif_msg.MessageID})
			
		// Отправки уведомления о успешном снятии антирежима
		bot.SendMessage(
			&telego.SendMessageParams{
				ChatID: update.Message.Chat.ChatID(),
				Text: "✅ Антирейд режим отключен.\nСвобода слова восстановлена!",
				ReplyParameters: &telego.ReplyParameters{
					MessageID: update.Message.MessageID,
				},
			})
		fmt.Printf(
			"Чат %v %v удалён из списка защиты\n", update.Message.Chat.ID, update.Message.Chat.Title)
	} else {
		fmt.Printf(
			"Чат %v %v не найден в списке защиты\n", update.Message.Chat.ID, update.Message.Chat.Title)
	}
}

func IsProtect(chatID int64) bool {
	if res, exitst := ProtectionChats[chatID]; exitst{
		return res.IsProtected
	}
	return false
}

func IsAutoban(chatID int64) bool {
	if res, exitst := ProtectionChats[chatID]; exitst{
		return res.Autoban
	}
	return false
}
