package card

type Comment struct {
	CommentId int    `json:"comment_id"`
	CardId    int    `json:"card_id"`
	Author    string `json:"author"`
	Body      string `json:"body"`
	Likes     int    `json:"likes"`
	Liked     bool   `json:"liked"`
}
