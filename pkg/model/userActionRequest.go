package model

type UserActionRequest struct {
	Action          string `json:"action"`
	UserCredentials `json:"user_credentials"`
	Target          string `json:"target"`
}
