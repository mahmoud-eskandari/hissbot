package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	//SendMessage SendMessage
	SendMessage = 1
	//ReplyTo ReplyTo
	ReplyTo = 2
	//FindSpecial FindSpecial
	FindSpecial = 3
)

//CheckCommand Command Dispatcher
func CheckCommand(update tgbotapi.Update, cmd string) (*tgbotapi.MessageConfig, *Command, error) {
	command := strings.Split(cmd, " ")
	switch command[0] {
	case "/start":
		return startDispatcher(&update, command)

	case "/reply":
		return replyMessage(&update)

	case "/block":
		return blockUser(&update)

	case "/unblock":
		return unblockUser(&update)

	case "/link":
		return sendLink(&update)

	case "/newmsg":
		return newMessage(&update)

	default:
		msg := tgbotapi.NewMessage(GetChatID(&update), IDontKnown)
		msg.ReplyMarkup = DefaultMenu
		return &msg, &Command{Command: 0}, nil
	}
}

// write new message
func newMessage(update *tgbotapi.Update) (*tgbotapi.MessageConfig, *Command, error) {
	// Get New Messages
	var msgs []Message
	err := Db.Find(&msgs, "seen=0 AND receiver=?", update.Message.From.ID).Error
	if err != nil || len(msgs) == 0 {
		return nil, nil, errors.New(YouHaveNotNewMessage)
	}
	for i := range msgs {
		_, err := Bot.Send(ShareMessage(&msgs[i], update.Message.Chat.ID, true))
		if err == nil {
			// Send Seen Message to sender
			seen := tgbotapi.NewMessage(int64(msgs[i].Sender), SeenMessage)
			seen.ReplyToMessageID = msgs[i].MessageID
			_, _ = Bot.Send(seen)
			msgs[i].Seen = 1
			err = Db.Save(&msgs[i]).Error
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	}

	return nil, &Command{Command: 0}, nil
}

func replyMessage(update *tgbotapi.Update) (*tgbotapi.MessageConfig, *Command, error) {
	message := Message{}
	err := Db.First(&message, "update_id=? AND receiver=?",
		strings.Replace(update.CallbackQuery.Data, "/reply ", "", 1),
		GetChatID(update)).Error

	if err != nil || message.UpdateID == 0 {
		return nil, nil, errors.New(MessageNotFound)
	}

	if checkForBlock(message.Sender, GetChatID(update)) {
		msg := tgbotapi.NewMessage(GetChatID(update), ForbiddenMessage)
		return &msg, &Command{}, nil
	}

	// Check Current user data
	_, _, _ = CurrentUser(update)
	reply := tgbotapi.NewMessage(GetChatID(update), fmt.Sprintf(SendingReply, makeUserHash(message.Sender, message.Receiver)))
	return &reply, &Command{Command: ReplyTo, CommandData: int64(message.UpdateID)}, nil
}

func sendMessage(update *tgbotapi.Update, id string) (*tgbotapi.MessageConfig, *Command, error) {
	user := User{}
	err := Db.First(&user, "link=?", id).Error
	if err != nil || user.ID == 0 {
		return nil, nil, errors.New(UserNotFound)
	}
	// check for block
	if checkForBlock(user.ID, GetChatID(update)) {
		msg := tgbotapi.NewMessage(GetChatID(update), ForbiddenMessage)
		return &msg, &Command{}, nil
	}
	// Check Current user data
	_, _, _ = CurrentUser(update)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(SendMessageStart, user.FirstName))
	msg.ReplyMarkup = CancelMenu
	return &msg, &Command{Command: SendMessage, CommandData: user.ID}, nil
}

func sendLink(update *tgbotapi.Update) (*tgbotapi.MessageConfig, *Command, error) {
	user, _, err := CurrentUser(update)
	if err != nil {
		return nil, nil, errors.New(ErrorOccured)
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(MyLink, user.FirstName, Bot.Self.UserName, user.Link))
	return &msg, &Command{Command: SendMessage, CommandData: user.ID}, nil
}

func startDispatcher(update *tgbotapi.Update, command []string) (*tgbotapi.MessageConfig, *Command, error) {
	if len(command) == 2 {
		return sendMessage(update, command[1])
	}

	msg := tgbotapi.NewMessage(GetChatID(update), IDontKnown)
	msg.ReplyMarkup = DefaultMenu
	return &msg, &Command{}, nil
}
