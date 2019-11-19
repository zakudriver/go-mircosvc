package common

import "time"

type (
	Event string
	Literal string
)

func (c Code) toString() string {
	return string(c)
}

func (c Code) Code() int32 {
	return int32(c)
}

const (
	Message          Event = "Message"
	SubscribeMessage Event = "SubscribeMessage"
	AlreadyMessage   Event = "AlreadyMessage"

	// CodeExpiration = 10 * 60

	SessionKey    = "SessionKey"
	AuthHeaderKey = "Authorization"
	ServerAuthKey = "Bearer"

	MaxAge = 24 * time.Hour / time.Second
)
