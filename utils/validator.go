package utils

import (
	"github.com/mymmrac/telego"
)

func CheckIsAdmin(bot *telego.Bot, meId int64, chatId int64) bool {
	admins, _ := bot.GetChatAdministrators(
		&telego.GetChatAdministratorsParams{ChatID: telego.ChatID{ID: chatId}},
	)
	for _, admin := range admins{
		if admin.MemberUser().ID == meId{return true}

	}
	return false
}