package card

type Card struct {
	BasicCard
	Comments []Comment
}

func (Card) GetType() string {
	return "Card"
}
