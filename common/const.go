package common

import "time"

type (
	Event string
	RoleLevel int8
)

const (
	// Message          Event = "Message"
	// SubscribeMessage Event = "SubscribeMessage"
	// AlreadyMessage   Event = "AlreadyMessage"

	// CodeExpiration = 10 * 60

	SessionKey = "SessionKey"
	CookieName = "Authorization"
	UIDKey     = "UID"
	RoleIDKey  = "RoleID"

	MaxAge = 24 * time.Hour / time.Second

	RootUser  RoleLevel = 0
	GuestUser RoleLevel = 1
)
