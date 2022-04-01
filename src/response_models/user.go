package response_models

type User struct {
	Id              string `json:"id"`
	Joined          string `json:"joined"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	Bio             string `json:"bio"`
	ProfilePhotoId  string `json:"profilePhotoId"`
	FollowersCount  string `json:"followersCount"`
	FollowingsCount string `json:"followingsCount"`
	PoemsCount      string `json:"poemsCount"`
	LikesCount      string `json:"likesCount"`
	CommentsCount   string `json:"commentsCount"`
	IsFollowing     string `json:"isFollowing"`
}

type UserMin struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	ProfilePhotoId string `json:"profilePhotoId"`
	IsFollowing    bool   `json:"isFollowing"`
}
