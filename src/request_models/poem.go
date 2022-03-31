package request_models

type PoemAddForm struct {
	AuthToken string   `json:"authToken" binding:"required"`
	UserId    string   `json:"userId" binding:"required"`
	Title     string   `json:"title" binding:"required"`
	Verses    []string `json:"verses" binding:"required"`
}

type PoemUpdateForm struct {
	AuthToken string   `json:"authToken" binding:"required"`
	UserId    string   `json:"userId" binding:"required"`
	PoemId    string   `json:"poemId" binding:"required"`
	Title     string   `json:"title" binding:"required"`
	Verses    []string `json:"verses" binding:"required"`
}

type PoemDeleteForm struct {
	AuthToken string `json:"authToken" binding:"required"`
	UserId    string `json:"userId" binding:"required"`
	PoemId    string `json:"poemId" binding:"required"`
}

type PoemLikeForm struct {
	AuthToken string `json:"authToken" binding:"required"`
	UserId    string `json:"userId" binding:"required"`
	PoemId    string `json:"poemId" binding:"required"`
}
