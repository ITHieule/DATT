package services

import (
	"errors"
	"fmt"
	"web-api/internal/api/until"
	"web-api/internal/pkg/database"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/types"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	*BaseService
}

var User = &UserService{}

func (s *UserService) UpdateUserSevice(requestParams *request.CreateUserRequest) ([]types.User, error) {
	var Sizes []types.User

	db, err := database.FashionBusiness()
	if err != nil {
		return nil, err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	var hashedPassword string
	if requestParams.Password != "" {
		hashedPassword, err = HashPassword(requestParams.Password)
		if err != nil {
			return nil, err
		}
	}

	query := `
        UPDATE users
        SET password_hash = COALESCE(?, password_hash),
		name = ?,
            email = ?,
            phone = ?,
			avatar_url = ?
        WHERE id = ?
    `

	err = db.Exec(query,
		hashedPassword,
		requestParams.Name,
		requestParams.Email,
		requestParams.Phone,
		requestParams.AvatarURL,
		requestParams.ID,
	).Error

	if err != nil {
		return nil, err
	}

	return Sizes, nil
}

func (s *UserService) Register(requestParams *request.CreateUserRequest) ([]types.User, error) {
	var User []types.User

	db, err := database.FashionBusiness()
	if err != nil {
		return nil, err
	}

	hashedPassword, err := HashPassword(requestParams.Password) // ✅ FE gửi "password", BE hash rồi lưu
	if err != nil {
		return nil, err
	}

	query := "INSERT INTO users (name, email, password_hash, phone) VALUES (?, ?, ?, ?)"
	err = db.Raw(query,
		requestParams.Name,
		requestParams.Email,
		hashedPassword,
		requestParams.Phone,
	).Scan(&User).Error

	if err != nil {
		fmt.Println("❌ Query execution error:", err)
	}

	return User, nil
}
func (s *UserService) Login(requestParams *request.LoginRequests) (string, error) {
	var user types.User

	db, err := database.FashionBusiness()
	if err != nil {
		return "", err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	query := "SELECT * FROM users WHERE email = ?"
	err = db.Raw(query, requestParams.Email).Scan(&user).Error

	if err != nil {
		return "", errors.New("email hoặc mật khẩu không đúng")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password_hash), []byte(requestParams.Password)) // ✅ FE gửi "password", BE kiểm tra hash
	if err != nil {
		return "", errors.New("email hoặc mật khẩu không đúng")
	}

	token, err := until.GenerateJWT(user.ID, user.Role, user.Name)
	if err != nil {
		return "", err
	}

	return token, nil
}

func HashPassword(Password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
