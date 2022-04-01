package response_models

type Poem struct {
	Id            string   `json:"id"`
	User          UserMin  `json:"user"`
	Title         string   `json:"title"`
	PublishedOn   string   `json:"publishedOn"`
	Verses        []string `json:"verses"`
	CommentsCount int      `json:"commentsCount"`
	LikesCount    int      `json:"likesCount"`
	IsLiked       bool     `json:"isLiked"`
}
