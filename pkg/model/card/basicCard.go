package card

type BasicCard struct {
	CardId      int    `json:"card_id"`
	Author      string `json:"author"`
	Image       string `json:"image"`
	Likes       int    `json:"likes"`
	Liked       bool   `json:"liked"`
	Description string `json:"description"`
}
