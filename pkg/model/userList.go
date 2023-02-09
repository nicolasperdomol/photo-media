package model

type UserList struct {
	Relation  string   `json:"relation"`
	Usernames []string `json:"user_list"`
}

func (UserList) GetType() string {
	return "UserList"
}
