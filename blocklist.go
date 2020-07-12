package main

import (
	"errors"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func blockUser(update *tgbotapi.Update) (*tgbotapi.MessageConfig, *Command, error) {
	blocked := DecodeID(strings.Replace(update.CallbackQuery.Data, "/block ", "", 1))
	suspend := Suspend{}
	suspend.OwnerID = GetChatID(update)
	suspend.BlockedUser = blocked
	err := Db.Create(&suspend).Error
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate") {
			return nil, nil, errors.New(DuplicateBlock)
		}
		return nil, nil, errors.New(ErrorOccured)
	}

	reply := tgbotapi.NewMessage(GetChatID(update), SuccessfulSuspend)
	reply.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏳️ آنبلاک کاربر", "/unblock "+EncodeID(blocked)),
		),
	)
	return &reply, &Command{}, nil
}

func unblockUser(update *tgbotapi.Update) (*tgbotapi.MessageConfig, *Command, error) {
	blocked := DecodeID(strings.Replace(update.CallbackQuery.Data, "/unblock ", "", 1))

	err := Db.Exec("DELETE FROM suspends WHERE owner_id=? AND blocked_user=?", GetChatID(update), blocked).Error
	if err != nil {
		return nil, nil, errors.New(ErrorOccured)
	}
	reply := tgbotapi.NewMessage(GetChatID(update), SuccessfulUnSuspend)
	return &reply, &Command{}, nil
}

//checkForBlock Check for blocked senders by receiver
func checkForBlock(owner int64, sender int64) bool {
	suspend := Suspend{}
	Db.First(&suspend, "owner_id=? AND blocked_user=?", owner, sender)
	return suspend.OwnerID > 0
}
