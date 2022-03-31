package request_models

type SignInForm struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SignUpForm struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type PasswordResetRequestForm struct {
	Email string `json:"email" binding:"required"`
}

type PasswordResetForm struct {
	Email      string `json:"email" binding:"required"`
	Password   string `json:"password" binding:"required"`
	ResetToken string `json:"resetToken" binding:"required"`
}
