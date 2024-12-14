package utils

import (
	"github.com/mymmrac/telego"
)

func BotIsAdmin(bot *telego.Bot, meId int64, chatId int64) bool {
	admins, _ := bot.GetChatAdministrators(
		&telego.GetChatAdministratorsParams{ChatID: telego.ChatID{ID: chatId}},
	)
	for _, admin := range admins{
		if admin.MemberUser().ID == meId{return true}
	}
	return false
}

func UserIsAdmin(bot *telego.Bot, update *telego.Update) bool {
	admins, _ := bot.GetChatAdministrators(
		&telego.GetChatAdministratorsParams{ChatID: telego.ChatID{ID: update.Message.Chat.ID}},
	)
	for _, admin := range admins{
		if admin.MemberUser().ID == update.Message.From.ID{return true}
	}
	return false
}