package request

type CreateUserRequest struct {
	ID         int    `json:"id"`
	Name       string `json:"name" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required"` // ✅ FE gửi "password", không gửi "password_hash"
	Phone      string `json:"phone"`
	AvatarURL  string `json:"avatar_url"`
	Created_at string `json:"created_at"`
	Role       string `json:"role"`
}

type LoginRequests struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password"`
}
