package main

import (
	"encoding/json"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type photo struct {
	FileID string `json:"file_id"`
}
type voice struct {
	FileID string `json:"file_id"`
}
type sticker struct {
	FileID string `json:"file_id"`
}

//ShareMessage Share message between users
func ShareMessage(message *Message, receiver int64, flags bool) tgbotapi.Chattable {
	uid := ""
	if flags {
		uid = fmt.Sprintf(UserText, makeUserHash(message.Sender, receiver))
	}

	// Photo
	var photos []photo
	err := json.Unmarshal([]byte(message.Photo), &photos)
	if err == nil && photos != nil && len(photos) >= 2 {
		msg := tgbotapi.NewPhotoShare(receiver, photos[1].FileID)
		msg.Caption = message.Caption + uid
		if flags {
			msg.ReplyMarkup = SetReplyReportKey(message.Sender, message.UpdateID)
		}
		return msg
	}

	// Stickers
	stickerMsg := sticker{}
	err = json.Unmarshal([]byte(message.Sticker), &stickerMsg)
	if err == nil && stickerMsg.FileID != "" {
		msg := tgbotapi.NewStickerShare(receiver, stickerMsg.FileID)
		if flags {
			msg.ReplyMarkup = SetReplyReportKey(message.Sender, message.UpdateID)
		}
		return msg
	}

	// Voice
	voiceMsg := voice{}
	err = json.Unmarshal([]byte(message.Voice), &voiceMsg)
	if err == nil && voiceMsg.FileID != "" {
		msg := tgbotapi.NewVoiceShare(receiver, voiceMsg.FileID)
		msg.Caption = message.Caption + uid
		if flags {
			msg.ReplyMarkup = SetReplyReportKey(message.Sender, message.UpdateID)
		}
		return msg
	}

	// Animation
	animationMsg := voice{}
	err = json.Unmarshal([]byte(message.Voice), &animationMsg)
	if err == nil && animationMsg.FileID != "" {
		msg := tgbotapi.NewAnimationShare(receiver, animationMsg.FileID)
		msg.Caption = message.Caption + uid
		if flags {
			msg.ReplyMarkup = SetReplyReportKey(message.Sender, message.UpdateID)
		}
		return msg
	}

	// File
	docMsg := voice{}
	err = json.Unmarshal([]byte(message.Document), &docMsg)
	if err == nil && docMsg.FileID != "" {
		msg := tgbotapi.NewDocumentShare(receiver, docMsg.FileID)
		msg.Caption = message.Caption + uid
		if flags {
			msg.ReplyMarkup = SetReplyReportKey(message.Sender, message.UpdateID)
		}
		return msg
	}

	// Audio
	audioMsg := voice{}
	err = json.Unmarshal([]byte(message.Audio), &audioMsg)
	if err == nil && audioMsg.FileID != "" {
		msg := tgbotapi.NewAudioShare(receiver, audioMsg.FileID)
		msg.Caption = message.Caption + uid
		if flags {
			msg.ReplyMarkup = SetReplyReportKey(message.Sender, message.UpdateID)
		}
		return msg
	}

	if len(message.Text) == 0 {
		message.Text = FormatNotSupported
	}

	msg := tgbotapi.NewMessage(receiver, message.Text+uid)
	if flags {
		msg.ReplyMarkup = SetReplyReportKey(message.Sender, message.UpdateID)
	}
	return msg
}

//CreateMessage Create message on database
func CreateMessage(update *tgbotapi.Update, receiver int64, reply int) bool {
	photo, _ := json.Marshal(update.Message.Photo)
	voice, _ := json.Marshal(update.Message.Voice)
	audio, _ := json.Marshal(update.Message.Audio)
	sticker, _ := json.Marshal(update.Message.Sticker)
	document, _ := json.Marshal(update.Message.Document)
	animation, _ := json.Marshal(update.Message.Animation)
	message := Message{
		UpdateID:  update.UpdateID,
		MessageID: update.Message.MessageID,
		ReplyTo:   reply,
		Text:      update.Message.Text + update.Message.Caption,
		Photo:     string(photo),
		Voice:     string(voice),
		Audio:     string(audio),
		Sticker:   string(sticker),
		Document:  string(document),
		Animation: string(animation),
		Caption:   update.Message.Caption,
		Sender:    update.Message.Chat.ID,
		Receiver:  receiver,
	}
	return Db.Create(&message).Error == nil
}
