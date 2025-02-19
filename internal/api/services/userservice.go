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
	if requestParams.Password != "" { // Kiểm tra nếu mật khẩu được cung cấp
		hashedPassword, err = HashPassword(requestParams.Password)
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

func (s *UserService) Register(requestParams *request.CreateUserRequest) ([]types.User, error) {
	var User []types.User

	// Kết nối cơ sở dữ liệu
	db, err := database.FashionBusiness()
	if err != nil {
		fmt.Println("❌ Database connection error:", err)
		return nil, err
	}

	// Log mật khẩu trước khi hash
	fmt.Println("🔍 Mật khẩu nhận từ FE:", requestParams.Password)

	// Mã hóa mật khẩu
	hashedPassword, err := HashPassword(requestParams.Password) // ✅ FE gửi "password", BE hash rồi lưu
	if err != nil {
		return nil, err
	}

	// Log mật khẩu đã hash
	fmt.Println("✅ Mật khẩu sau khi hash:", hashedPassword)

	// Lưu vào database với cột "password_hash"
	query := "INSERT INTO users (name, email, password_hash, phone) VALUES (?, ?, ?, ?)"
	err = db.Raw(query,
		requestParams.Name,
		requestParams.Email,
		hashedPassword, // ✅ Lưu vào cột "password_hash" trong DB
		requestParams.Phone,
	).Scan(&User).Error

	if err != nil {
		fmt.Println("❌ Query execution error:", err)
	}

	return User, nil
}
func (s *UserService) Login(requestParams *request.LoginRequests) (string, error) {
	var user types.User

	// Log dữ liệu nhận từ FE
	fmt.Println("📥 Nhận request login:", requestParams)

	// Kết nối cơ sở dữ liệu
	db, err := database.FashionBusiness()
	if err != nil {
		fmt.Println("❌ Lỗi kết nối DB:", err)
		return "", err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Truy vấn thông tin người dùng dựa trên email
	query := "SELECT * FROM users WHERE email = ?"
	err = db.Raw(query, requestParams.Email).Scan(&user).Error

	if err != nil {
		fmt.Println("❌ Không tìm thấy user với email:", requestParams.Email)
		return "", errors.New("email hoặc mật khẩu không đúng")
	}

	// Log dữ liệu user từ DB
	fmt.Println("✅ Tìm thấy user trong DB:", user)

	// So sánh mật khẩu đã hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password_hash), []byte(requestParams.Password)) // ✅ FE gửi "password", BE kiểm tra hash
	if err != nil {
		fmt.Println("❌ Mật khẩu không khớp cho user:", requestParams.Email)
		fmt.Println("🔍 Mật khẩu nhập vào:", requestParams.Password)
		fmt.Println("🔍 Mật khẩu mã hóa trong DB:", user.Password_hash)
		return "", errors.New("email hoặc mật khẩu không đúng")
	}

	// Tạo JWT token
	token, err := until.GenerateJWT(user.ID, user.Role, user.Name)
	if err != nil {
		fmt.Println("❌ Lỗi khi tạo JWT token:", err)
		return "", err
	}

	// Log token trước khi trả về
	fmt.Println("✅ Đăng nhập thành công, token:", token)

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
