package main

import (
	"fmt"
	"strconv"
)

// MessageEntity contains information about data in a Message.
type MessageEntity struct {
	Type   string `json:"type"`
	Offset int    `json:"offset"`
	Length int    `json:"length"`
	URL    string `json:"url"`  // optional
	User   *User  `json:"user"` // optional
}

// Message is a message recieved from telegram
type Message struct {
	MessageID             int              `json:"message_id" binding:"required"`
	From                  *User            `json:"from"` // optional
	Date                  int              `json:"date"`
	Chat                  *Chat            `json:"chat"`
	ForwardFrom           *User            `json:"forward_from"`            // optional
	ForwardFromChat       *Chat            `json:"forward_from_chat"`       // optional
	ForwardFromMessageID  int              `json:"forward_from_message_id"` // optional
	ForwardDate           int              `json:"forward_date"`            // optional
	ReplyToMessage        *Message         `json:"reply_to_message"`        // optional
	EditDate              int              `json:"edit_date"`               // optional
	Text                  string           `json:"text"`                    // optional
	Entities              *[]MessageEntity `json:"entities"`                // optional
	Document              *Document        `json:"document"`                // optional
	Photo                 *[]PhotoSize     `json:"photo"`                   // optional
	Sticker               *Sticker         `json:"sticker"`                 // optional
	Caption               string           `json:"caption"`                 // optional
	Contact               *Contact         `json:"contact"`                 // optional
	Location              *Location        `json:"location"`                // optional
	Venue                 *Venue           `json:"venue"`                   // optional
	NewChatMembers        *[]User          `json:"new_chat_members"`        // optional
	LeftChatMember        *User            `json:"left_chat_member"`        // optional
	NewChatTitle          string           `json:"new_chat_title"`          // optional
	DeleteChatPhoto       bool             `json:"delete_chat_photo"`       // optional
	GroupChatCreated      bool             `json:"group_chat_created"`      // optional
	SuperGroupChatCreated bool             `json:"supergroup_chat_created"` // optional
	ChannelChatCreated    bool             `json:"channel_chat_created"`    // optional
	MigrateToChatID       int64            `json:"migrate_to_chat_id"`      // optional
	MigrateFromChatID     int64            `json:"migrate_from_chat_id"`    // optional
	PinnedMessage         *Message         `json:"pinned_message"`          // optional
}

func (message Message) ToString() string {

	id := strconv.FormatInt(message.Chat.ID, 10)

	chat := message.Chat.Title + "@" + id
	if message.Chat.UserName != "" {
		chat = message.Chat.UserName + "@" + id
	}

	content := message.Text
	if message.Sticker != nil {
		content = "<sticker>"
	} else if message.Photo != nil {
		content = "<photo>"
	}

	displayname := message.From.FirstName + " " + message.From.LastName

	return fmt.Sprintf("[%s] %s: %s", chat, displayname, content)
}
