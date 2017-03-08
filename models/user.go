package models

import "time"

type User struct {
	Id          string
	Username    string
	Email       string
	Token       string
	SnipCount   int64
	RunCount    int64
	ClonedCount int64
	Created     time.Time
}

var Principal *User

type ChangePasswordParams struct {
	Username         string
	Email            string
	ResetToken       string
	ExistingPassword string
	NewPassword      string
}
