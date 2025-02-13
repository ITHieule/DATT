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

func (s *UserService) Register(requestParams *request.CreateUserRequest) ([]types.User, error) {
	var User []types.User
	// kết nối cơ sở dữ liệu
	db, err := database.FashionBusiness()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return nil, err
	}
	// Mã hóa mật khẩu
	hashedPassword, err := HashPassword(requestParams.Password_hash)
	if err != nil {

		return nil, err
	}
	requestParams.Password_hash = hashedPassword
	// Kiểm tra dữ liệu đầu vào
	dbInstance, _ := db.DB()
	defer dbInstance.Close()
	query := "INSERT INTO users (name, email, password_hash,phone) VALUES ( ?, ?, ?, ?)"
	err = db.Raw(query,
		requestParams.Name,
		requestParams.Email,
		requestParams.Password_hash,
		requestParams.Phone,
	).Scan(&User).Error
	if err != nil {
		fmt.Println("Query execution error:", err)
	}

	return User, nil
}

func (s *UserService) UpdateUserSevice(requestParams *request.CreateUserRequest) ([]types.User, error) {
	var Sizes []types.User

	// Kết nối database
	db, err := database.FashionBusiness()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return nil, err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Mã hóa mật khẩu nếu có thay đổi
	var hashedPassword string
	if requestParams.Password_hash != "" { // Kiểm tra nếu mật khẩu được cung cấp
		hashedPassword, err = HashPassword(requestParams.Password_hash)
		if err != nil {
			fmt.Println("Password hashing error:", err)
			return nil, err
		}
	}

	// Truy vấn SQL cập nhật thông tin người dùng
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
		hashedPassword, // Mật khẩu đã mã hóa hoặc NULL để giữ nguyên
		requestParams.Name,
		requestParams.Email,
		requestParams.Phone,
		requestParams.AvatarURL,
		requestParams.ID,
	).Error

	if err != nil {
		fmt.Println("Query execution error:", err)
		return nil, err
	}

	return Sizes, nil
}

func (s *UserService) Login(requestParams *request.LoginRequests) (string, error) {
	var user types.User

	// Kết nối cơ sở dữ liệu
	db, err := database.FashionBusiness()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return "", err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Truy vấn thông tin người dùng dựa trên email
	query := "SELECT * FROM users WHERE email = ?"
	err = db.Raw(query, requestParams.Email).Scan(&user).Error

	if err != nil {
		fmt.Println("Query execution error:", err)
		return "", err
	}

	// So sánh mật khẩu đã mã hóa với mật khẩu người dùng nhập vào
	err = bcrypt.CompareHashAndPassword([]byte(user.Password_hash), []byte(requestParams.Password_hash))
	if err != nil {
		fmt.Println("Password mismatch for user:", requestParams.Email)
		return "", errors.New("invalid credentials") // Trả về lỗi nếu mật khẩu không khớp
	}
	// Tạo JWT token
	token, err := until.GenerateJWT(user.ID, user.Role, user.Name)
	if err != nil {
		fmt.Println("Error generating token:", err)
		return token, nil
	}

	// Trả về thông tin người dùng và token
	return token, nil

}

// Hàm mã hóa mật khẩu
func HashPassword(Password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
