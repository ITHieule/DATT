package until

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Khóa bí mật để ký token
var jwtSecret = []byte("your_secret_key")

// Hàm tạo JWT token
func GenerateJWT(userID int, Role string, Username string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"role":     Role,
		"username": Username,
		"exp":      time.Now().Add(time.Hour * 2).Unix(), // Hạn token 2 giờ
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
