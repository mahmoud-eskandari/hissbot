package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//TextCommands Check for menu's text commands
func TextCommands(update *tgbotapi.Update) bool {
	switch update.Message.Text {
	case Cancel:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, OkWhatIDo)
		msg.ReplyMarkup = DefaultMenu
		_, _ = Bot.Send(msg)
		// Reset Command State
		UserCommandState[update.Message.From.ID] = Command{0, 0}
		return true
	case MyLinkCommand:
		msg, _, _ := sendLink(update)
		msg.ReplyMarkup = DefaultMenu
		_, _ = Bot.Send(msg)
		return true
	case IncreaseAwards:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, Awards)
		msg.ReplyMarkup = DefaultMenu
		_, _ = Bot.Send(msg)
		return true
	case HelpCommand:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, HelpMessage)
		msg.ReplyMarkup = DefaultMenu
		_, _ = Bot.Send(msg)
		return true
	case ConnectToSpecific:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, SpecificMessage)
		UserCommandState[update.Message.From.ID] = Command{Command: 3, CommandData: 1}
		msg.ReplyMarkup = CancelMenu
		_, _ = Bot.Send(msg)
		return true
	case ConnectToAny:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, ToAnyMessage)
		msg.ReplyMarkup = DefaultMenu
		_, _ = Bot.Send(msg)
		return true
	default:
		return false
	}
}

//IsInlineCommand check string for inline command
func IsInlineCommand(str string) bool {
	switch str {
	case Cancel:
		return true
	case ConnectToAny:
		return true
	case ConnectToSpecific:
		return true
	case AnonymousToGroup:
		return true
	case MyLinkCommand:
		return true
	case IncreaseAwards:
		return true
	case HelpCommand:
		return true
	}
	return false
}
