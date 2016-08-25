package system

type User struct {
	Id      int64 `json:"id"`
	Username string `json:"username"`
	Email 	string `json:"email"`
	Host    string `json:"host"`
	Token    string `json:"token"`
}

func (u *User) Err() string {
	if len(u.Token) < 1 { return "Failed to authenticate, bad username or password."}
	return ""
}
