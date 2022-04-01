package response_models

type Comment struct {
	Id           string  `json:"id"`
	User         UserMin `json:"user"`
	CreatedOn    string  `json:"createdOn"`
	Text         string  `json:"text"`
	PoemId       string  `json:"poemId"`
	RepliesCount int     `json:"repliesCount"`
	ReplyTo      string  `json:"replyTo"`
}
