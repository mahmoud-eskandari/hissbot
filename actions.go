package main

import (
	"errors"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//DoCommand Do actions
func DoCommand(update tgbotapi.Update, command Command) (*tgbotapi.MessageConfig, error) {
	if update.Message != nil && update.Message.Text == Cancel {
		msg := tgbotapi.NewMessage(GetChatID(&update), Canceled)
		msg.ReplyMarkup = DefaultMenu
		return &msg, nil
	}
	//TODO: move cases to diffrent functions if need more params
	switch command.Command {
	case SendMessage:
		if command.CommandData == update.Message.Chat.ID {
			return nil, errors.New(YouCantSendToYourSelf)
		}

		if checkForBlock(command.CommandData, GetChatID(&update)) {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, ForbiddenMessage)
			msg.ReplyMarkup = DefaultMenu
			return &msg, nil
		}

		if !CreateMessage(&update, command.CommandData, 0) {
			// Send Error
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, ErrorOccured)
			msg.ReplyMarkup = DefaultMenu
			return &msg, nil
		}

		// Notice Sender
		msg := tgbotapi.NewMessage(GetChatID(&update), SentSuccessfully)
		msg.ReplyMarkup = DefaultMenu
		_, _ = Bot.Send(msg)

		// Notice Receiver
		msg = tgbotapi.NewMessage(command.CommandData, NewMessage)
		return &msg, nil

	case ReplyTo:
		prevMessage := Message{}
		Db.First(&prevMessage, "update_id=?", command.CommandData)
		if checkForBlock(prevMessage.Sender, GetChatID(&update)) {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, ForbiddenMessage)
			msg.ReplyMarkup = DefaultMenu
			return &msg, nil
		}
		if !CreateMessage(&update, prevMessage.Sender, int(command.CommandData)) {
			// Send Error
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, ErrorOccured)
			msg.ReplyMarkup = DefaultMenu
			return &msg, nil
		}
		// Notice Sender
		msg := tgbotapi.NewMessage(prevMessage.Receiver, SentSuccessfully)
		msg.ReplyMarkup = DefaultMenu
		_, _ = Bot.Send(msg)
		// Notice Receiver
		msg = tgbotapi.NewMessage(prevMessage.Sender, NewMessage)
		return &msg, nil

	case FindSpecial:
		usr := User{}

		if strings.HasPrefix(update.Message.Text, "@") {
			Db.First(&usr, "user_name=?", strings.Replace(update.Message.Text, "@", "", -1))
		} else {
			if update.Message != nil && update.Message.ForwardFrom != nil {
				Db.First(&usr, "id=?", update.Message.ForwardFrom.ID)
			}
		}
		if usr.ID > 0 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(SpecialFind, Bot.Self.UserName, usr.Link))
			msg.ReplyMarkup = DefaultMenu
			return &msg, nil
		}
		//Reply
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, SpecialNotFound)
		msg.ReplyMarkup = DefaultMenu
		return &msg, nil

	}

	return nil, errors.New(IDontKnown)
}
