package request_models

type CommentAddForm struct {
	AuthToken string `json:"authToken" binding:"required"`
	UserId    string `json:"userId" binding:"required"`
	PoemId    string `json:"poemId" binding:"required"`
	Text      string `json:"text" binding:"required"`
	ReplyTo   string `json:"replyTo" binding:"-"`
}

type CommentDeleteForm struct {
	AuthToken string `json:"authToken" binding:"required"`
	UserId    string `json:"userId" binding:"required"`
	CommentId string `json:"commentId" binding:"required"`
}
