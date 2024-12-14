package utils

import (
	"fmt"
	"github.com/mymmrac/telego"
	"time"
)


var ProtectionChats map[int64]bool

func init() {
	// Инициализация глобальной переменной ProtectionChats
	ProtectionChats = make(map[int64]bool)
}


func AddChat(bot *telego.Bot, update *telego.Update) {
	ProtectionChats[update.Message.Chat.ID] = true
	notif_msg, _ := bot.SendMessage(
		&telego.SendMessageParams{
			ChatID: update.Message.Chat.ChatID(),
			Text: "🚨 ВКЛЮЧЕН АНТИРЕЙД РЕЖИМ 🚨\nВСЕ СООБЩЕНИЯ В ТЕЧЕНИИ 60 СЕКУНД БУДУТ УДАЛЯТЬСЯ"})
	fmt.Printf("Чат %v добавлен в защиту\n", update.Message.Chat.ID)
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
				Text: "✅ Антирейд режим отключен.\nТеперь можно свободно писать сообщения."})
		fmt.Printf("Чат %v удалён из списка защиты\n", update.Message.Chat.ID)
	} else {
		fmt.Printf("Чат %v не найден в списке защиты\n", update.Message.Chat.ID)
	}
}

func IsProtect(chatID int64) bool {
	return ProtectionChats[chatID]
}

