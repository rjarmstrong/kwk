package models

type User struct {
	Id      int64
	Username string
	Email 	string
	Token    string
	AliasCount int64
	RunCount int64
}