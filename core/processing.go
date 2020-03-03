package core

import (
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
)

// originalMessage is an unmarshaled raw message
type unmarshaledMessage struct {
	OriginTimestamp int    `json:"origin_server_ts"`
	Origin          string `json:"origin"`
	Sender          string `json:"sender"`
	EventID         string `json:"event_id"`
	RoomID          string `json:"room_id"`
	Content         struct {
		Message string `json:"body"`
		URL     string `json:"url"`
	} `json:"content"`
}

// Message is a processed message
type Message struct {
	Timestamp       int    `json:"ts"`
	OriginTimestamp int    `json:"origin_ts"`
	Origin          string `json:"origin"`
	Sender          string `json:"sender"`
	EventID         string `json:"event_id"`
	RoomID          string `json:"room_id"`
	Message         string `json:"message"`
	URL             string `json:"url"`
	ChatType        string `json:"chat_type"`
	Nick            string `json:"nick"`
}

// ProcessMessage returns a Message object from a raw message
func ProcessMessage(timestamp string, rawMessage string) (message Message, err error) {
	var unmarshaledMessage unmarshaledMessage
	if err = json.Unmarshal([]byte(rawMessage), &unmarshaledMessage); err != nil {
		return
	}
	timestampInt, err := strconv.Atoi(timestamp)
	if err != nil {
		return
	}
	nick, err := SenderToNick(unmarshaledMessage.Sender)
	if err != nil {
		return
	}

	message = Message{
		Timestamp:       timestampInt,
		OriginTimestamp: unmarshaledMessage.OriginTimestamp,
		Origin:          unmarshaledMessage.Origin,
		Sender:          unmarshaledMessage.Sender,
		EventID:         unmarshaledMessage.EventID,
		RoomID:          unmarshaledMessage.RoomID,
		Message:         unmarshaledMessage.Content.Message,
		URL:             unmarshaledMessage.Content.URL,
		ChatType:        "matrix",
		Nick:            nick,
	}
	return
}

// SenderToNick returns a short nickname from a Matrix full sender
func SenderToNick(sender string) (nick string, err error) {
	r := regexp.MustCompile("^@(.*):.*$")
	nick = r.ReplaceAllString(sender, "$1")
	if nick == "" {
		err = errors.New("cannot convert Matrix sender to a proper nickname")
	}
	return
}
