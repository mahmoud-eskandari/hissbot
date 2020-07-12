package main

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var DefaultMenu = tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{
	{
		Text:            ConnectToAny,
		RequestContact:  false,
		RequestLocation: false,
	},
}, []tgbotapi.KeyboardButton{
	{
		Text:            MyLinkCommand,
		RequestContact:  false,
		RequestLocation: false,
	},
	{
		Text:            ConnectToSpecific,
		RequestContact:  false,
		RequestLocation: false,
	},
}, []tgbotapi.KeyboardButton{
	{
		Text:            HelpCommand,
		RequestContact:  false,
		RequestLocation: false,
	},
	{
		Text:            Cancel,
		RequestContact:  false,
		RequestLocation: false,
	},
})

var CancelMenu = tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{
	{
		Text:            Cancel,
		RequestContact:  false,
		RequestLocation: false,
	},
})

func SetReplyReportKey(senderID int64, messageID int) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⛔️ بلاک", "/block "+EncodeID(senderID)),
			tgbotapi.NewInlineKeyboardButtonData("✍️ پاسخ", fmt.Sprintf("/reply %d", messageID)),
		),
	)
}
