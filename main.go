package main

import (
	"errors"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strings"
)

var Bot *tgbotapi.BotAPI
var UserCommandState map[int]Command

func main() {
	// Save last command state in this map
	UserCommandState = make(map[int]Command)
	// initialize config
	initConnections()
	defer CloseDB()
	// Migrate()
	// Start Bot
	err := errors.New("")
	Bot, err = tgbotapi.NewBotAPI(Str("TELEGRAM_API", ""))
	if err != nil {
		log.Panic(err)
	}
	// Start Debug
	Bot.Debug = Bool("DEBUG", false)
	if Bot.Debug {
		log.Printf("Authorized on account %s", Bot.Self.UserName)
	}

	// Fetch New Updates
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := Bot.GetUpdatesChan(u)
	for update := range updates {

		chatID := int(GetChatID(&update))
		if chatID == 0 {
			continue
		}
		// ignore bots
		if update.Message != nil && update.Message.From.IsBot {
			continue
		}

		command := ""
		// Integrate Inline and Query Commands
		if update.Message == nil {
			if update.CallbackQuery != nil {
				command = update.CallbackQuery.Data
			}
		} else {
			if strings.HasPrefix(update.Message.Text, "/") {
				command = update.Message.Text
				UserCommandState[update.Message.From.ID] = Command{}
			} else if IsInlineCommand(update.Message.Text) {
				UserCommandState[update.Message.From.ID] = Command{}
			}
		}
		// check for start UserCommandState
		if len(command) > 0 {
			msg, command, err := CheckCommand(update, command)
			if err != nil {
				_, _ = Bot.Send(tgbotapi.NewMessage(GetChatID(&update), err.Error()))
			} else {
				if command != nil {
					UserCommandState[chatID] = *command
				} else {
					UserCommandState[chatID] = Command{}
				}

				if msg != nil {
					_, _ = Bot.Send(msg)
				}
			}
			continue
		}

		// Check if is a text command
		if update.Message != nil && TextCommands(&update) {
			continue
		}

		// Do UserCommandState
		if command, ok := UserCommandState[chatID]; ok {
			if command.CommandData > 0 {
				msg, err := DoCommand(update, command)
				if err != nil {
					_, _ = Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
				} else {
					if msg != nil {
						_, _ = Bot.Send(msg)
					}
				}
				UserCommandState[chatID] = Command{Command: 0}
				continue
			}
		}

		msg := tgbotapi.NewMessage(int64(chatID), IDontKnown)
		msg.ReplyMarkup = DefaultMenu
		_, _ = Bot.Send(msg)
	}
}
