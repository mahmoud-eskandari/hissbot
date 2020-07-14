package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/labstack/gommon/random"
)

//Command Command
type Command struct {
	Command     int
	CommandData int64
}

//User User
type User struct {
	ID           int64     `gorm:"primary_key" json:"id"`
	Link         string    `gorm:"size:32" json:"link"`
	FirstName    string    `gorm:"size:255" json:"first_name"`
	LastName     string    `gorm:"size:255" json:"last_name"`     // optional
	UserName     string    `gorm:"size:255" json:"username"`      // optional
	LanguageCode string    `gorm:"size:255" json:"language_code"` // optional
	BlockedBot   bool      `json:"blocked_bot"`                   // optional
	IsBot        bool      `json:"is_bot"`                        // optional
	CreatedAt    time.Time `json:"created_at"`
}

//Suspend Suspend
type Suspend struct {
	ID          uint
	OwnerID     int64 `gorm:"primary_key;auto_increment:false" json:"owner_id"`
	BlockedUser int64 `gorm:"primary_key;auto_increment:false" json:"blocked_user"`
}

//Message Message
type Message struct {
	UpdateID  int    `gorm:"primary_key;auto_increment:false" json:"update_id"`
	MessageID int    `json:"message_id"`
	Text      string `json:"text"`
	ReplyTo   int    `json:"reply_to"`
	Photo     string `json:"photo" gorm:"type:text"`
	Voice     string `json:"voice" gorm:"type:text"`
	Audio     string `json:"audio" gorm:"type:text"`
	Sticker   string `json:"sticker" gorm:"type:text"`
	Document  string `json:"document" gorm:"type:text"`
	Animation string `json:"animation" gorm:"type:text"`
	Caption   string `json:"caption"`
	Sender    int64  `json:"sender"`
	Receiver  int64  `json:"receiver"`
	Seen      int    `json:"seen"`
}

//Migrate Simple migrate withoiut extrnal tools
func migrate() {
	Db.AutoMigrate(&User{}, &Message{}, &Suspend{})
}

//GetChatID get ChatID form message
func GetChatID(update *tgbotapi.Update) int64 {
	if update.Message == nil {
		if update.CallbackQuery == nil {
			return 0
		}
		return update.CallbackQuery.Message.Chat.ID
	}
	return update.Message.Chat.ID
}

//CurrentUser get user by session
func CurrentUser(update *tgbotapi.Update) (user User, new bool, err error) {
	user.ID = GetChatID(update)
	// First Check Redis
	if user.GetRedis() == nil {
		return
	}
	err = Db.First(&user, "id=?", GetChatID(update)).Error
	if err != nil || user.ID == 0 {
		new = true
		user = User{
			ID:           GetChatID(update),
			Link:         random.String(10),
			FirstName:    update.Message.From.FirstName,
			LastName:     update.Message.From.LastName,
			UserName:     update.Message.From.UserName,
			LanguageCode: update.Message.From.LanguageCode,
			BlockedBot:   false,
			IsBot:        update.Message.From.IsBot,
		}
		err = Db.Create(&user).Error
		if err != nil {
			return
		}
	}
	_ = SaveToRedis(fmt.Sprintf("user_%d", user.ID), user, time.Hour*6)
	return
}

//GetRedis fetch user from redis
func (usr *User) GetRedis() error {
	if usr.ID == 0 {
		return errors.New("ID is Empty")
	}
	s := Redis.Get(fmt.Sprintf("user_%d", usr.ID))
	if s.Err() != nil {
		return s.Err()
	}
	b, err := s.Bytes()
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, usr)
	if err != nil {
		return err
	}
	return nil
}
