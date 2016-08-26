package system

type User struct {
	Id      int64 `json:"id"`
	Username string `json:"username"`
	Email 	string `json:"email"`
	Host    string `json:"host"`
	Token    string `json:"token"`
	Message    string `json:"message"`
	Error    string `json:"error"`
}

func (u *User) Err() string {
	if len(u.Message)>0 || len(u.Error) >0 {
		return u.Message + " " + u.Error
	} else {
		if len(u.Token) < 1 {
			return "Failed to authenticate, bad username or password."
		}
	}
	return ""
}
