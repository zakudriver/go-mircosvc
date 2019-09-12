package common

import "time"

type (
	Code int32
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
	OK    Code = 0
	Error Code = 1

	Message          Event = "Message"
	SubscribeMessage Event = "SubscribeMessage"
	AlreadyMessage   Event = "AlreadyMessage"

	// CodeExpiration = 10 * 60

	SessionKey    = "SessionKey"
	AuthHeaderKey = "Authorization"
	ServerAuthKey = "Bearer"

	MaxAge = 24 * time.Hour / time.Second
)
