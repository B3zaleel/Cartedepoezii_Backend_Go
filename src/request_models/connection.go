package request_models

type ConnectionForm struct {
	AuthToken string `json:"authToken" binding:"required"`
	UserId    string `json:"userId" binding:"required"`
	FollowId  string `json:"followId" binding:"required"`
}
