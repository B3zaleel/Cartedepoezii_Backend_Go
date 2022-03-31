package request_models

type UserDeleteForm struct {
	AuthToken string `json:"authToken" binding:"required"`
	UserId    string `json:"userId" binding:"required"`
}

type UserUpdateForm struct {
	AuthToken          string `json:"authToken" binding:"required"`
	UserId             string `json:"userId" binding:"required"`
	Name               string `json:"name" binding:"required"`
	ProfilePhoto       string `json:"profilePhoto" binding:"-"`
	ProfilePhotoId     string `json:"profilePhotoId" binding:"required"`
	RemoveProfilePhoto bool   `json:"removeProfilePhoto" binding:"required"`
	Email              string `json:"email" binding:"required"`
	Bio                string `json:"bio" binding:"required"`
}
