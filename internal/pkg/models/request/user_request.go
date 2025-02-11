package request

type CreateUserRequest struct {
	Name      string `json:"name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
	Phone     string `json:"phone"`
	AvatarURL string `json:"avatar_url"`
}

type LoginRequests struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
