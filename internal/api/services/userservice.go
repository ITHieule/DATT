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

func (s *UserService) GetUserService() ([]types.User, error) {
	var records []types.User

	// Kết nối database
	db, err := database.FashionBusiness()
	if err != nil {
		return nil, fmt.Errorf("lỗi kết nối database: %w", err)
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Thực hiện truy vấn
	query := `SELECT * FROM users`
	err = db.Raw(query).Scan(&records).Error
	if err != nil {
		return nil, fmt.Errorf("lỗi truy vấn dữ liệu: %w", err)
	}

	return records, nil
}

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

func (s *UserService) Register(requestParams *request.CreateUserRequest) (*types.User, error) {
	db, err := database.FashionBusiness()
	if err != nil {
		return nil, errors.New("Lỗi kết nối cơ sở dữ liệu")
	}

	// Kiểm tra email đã tồn tại chưa
	var existingEmail types.User
	err = db.Raw("SELECT id FROM users WHERE email = ? LIMIT 1", requestParams.Email).Scan(&existingEmail).Error
	if err != nil {
		return nil, err
	}
	if existingEmail.ID != 0 {
		return nil, errors.New("email đã tồn tại")
	}

	// Kiểm tra số điện thoại đã tồn tại chưa
	var existingPhone types.User
	err = db.Raw("SELECT id FROM users WHERE phone = ? LIMIT 1", requestParams.Phone).Scan(&existingPhone).Error
	if err != nil {
		return nil, err
	}
	if existingPhone.ID != 0 {
		return nil, errors.New("số điện thoại đã tồn tại")
	}

	// Mã hóa mật khẩu
	hashedPassword, err := HashPassword(requestParams.Password)
	if err != nil {
		return nil, errors.New("Lỗi khi mã hóa mật khẩu")
	}

	// Tạo user mới
	var newUser types.User
	err = db.Raw("INSERT INTO users (name, email, password_hash, phone) VALUES (?, ?, ?, ?) RETURNING *",
		requestParams.Name, requestParams.Email, hashedPassword, requestParams.Phone).Scan(&newUser).Error

	if err != nil {
		return nil, errors.New("Lỗi khi tạo tài khoản")
	}

	return &newUser, nil
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
