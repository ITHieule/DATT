package services

import (
	"fmt"
	"os"
	"time"
	"web-api/internal/pkg/database"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/types"

	"github.com/dgrijalva/jwt-go"
)

// FILE - ADMIN SERVICE
type AdminService struct {
	*BaseService
}

var Admin = &AdminService{}

// func - LoginService

func (s *AdminService) LoginService(requestParams *request.LoginRequest) (string, error) {
	var admin types.Admin

	// Kết nối cơ sở dữ liệu
	db, err := database.RPConnection()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return "", err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Truy vấn lấy thông tin admin
	query := "SELECT adminid, username, password FROM report_admins WHERE username = ?"
	err = db.Raw(query, requestParams.Username).Scan(&admin).Error
	if err != nil {
		fmt.Println("Query error:", err)
		return "", err
	}

	// So sánh mật khẩu
	if admin.Password != requestParams.Password {
		return "", fmt.Errorf("invalid username or password")
	}

	// Tạo JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"adminid":  admin.Adminid,
		"username": admin.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token hết hạn sau 24 giờ
	})

	// Ký token
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		fmt.Println("Token signing error:", err)
		return "", err
	}

	return tokenString, nil
}

// func - GetReportByIsAnsweredAndRpidService
func (s *AdminService) GetReportByIsAnsweredAndRpidService(requestParams *request.GetReportByIsAnsweredAndRpidRequest) ([]types.ReportMng, error) {
	var rp_mng []types.ReportMng

	// Kết nối cơ sở dữ liệu
	db, err := database.RPConnection()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return nil, err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Truy vấn lấy dữ liệu
	query := `
	SELECT a.*, b.title as depname FROM report_mng a
	LEFT JOIN report_departments b
	ON b.dpid = a.dpid 
	WHERE rptid = ? 
	AND is_answered = ? 
	AND is_completed = ? 
	ORDER BY updatedate DESC
	`

	err = db.Raw(query, requestParams.Rptid, requestParams.IsAnswered, requestParams.IsCompleted).Scan(&rp_mng).Error
	if err != nil {
		fmt.Println("Query error:", err)
		return nil, err
	}

	return rp_mng, nil
}

// func - SendAnswerService
func (s *AdminService) SendAnswerService(requestParams *request.SendAnswerRequest) (string, error) {
	db, err := database.RPConnection()
	if err != nil {
		return "error", fmt.Errorf("database connection error: %v", err)
	}
	dbInstance, _ := db.DB()

	// Start a transaction
	tx := db.Begin()
	if tx.Error != nil {
		return "error", fmt.Errorf("failed to begin transaction: %v", tx.Error)
	}

	answerTime := time.Now().Format("2006-01-02 15:04:05")
	// Update reports answer
	query := fmt.Sprintf(`UPDATE reports SET answer = '%s', answer_time = '%s' where rpcode = '%s' and id='%d'`, requestParams.Answer, answerTime, requestParams.Rpcode, requestParams.Id)

	err = tx.Exec(query).Error
	if err != nil {
		tx.Rollback()
		return "error", err
	}

	// Update report_mng
	query = fmt.Sprintf(`UPDATE report_mng SET is_answered = 1, updatedate = '%s' where rpcode = '%s'`, answerTime, requestParams.Rpcode)
	err = tx.Exec(query).Error

	if err != nil {
		tx.Rollback()
		return "error", err
	}

	// Commit the transaction
	err = tx.Commit().Error
	if err != nil {
		return "error", err
	}

	dbInstance.Close()

	return "success", nil
}

// func - UpdateAnswerService
func (s *AdminService) UpdateAnswerService(requestParams *request.UpdateAnswerRequest) (string, error) {
	db, err := database.RPConnection()
	if err != nil {
		return "error", fmt.Errorf("database connection error: %v", err)
	}
	dbInstance, _ := db.DB()

	query := fmt.Sprintf(`UPDATE reports SET answer = '%s' WHERE id = '%d' AND rpcode = '%s'`, requestParams.Answer, requestParams.Id, requestParams.Rpcode)

	err = db.Exec(query).Error
	if err != nil {
		return "error", err
	}

	dbInstance.Close()

	return "success", nil
}

// func - GetRPTypesService
func (s *AdminService) GetRPTypesService() ([]types.ReportType, error) {
	var reportTypes []types.ReportType
	db, err := database.RPConnection()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return nil, err
	}
	dbInstance, _ := db.DB()
	// truy vấn lấy dữ liệu
	query := `SELECT * FROM report_types`

	err = db.Raw(query).Scan(&reportTypes).Error
	if err != nil {
		fmt.Println("Query error:", err)
		return nil, err
	}
	dbInstance.Close()
	return reportTypes, nil

}

// func - UpdateRPTypeService
func (s *AdminService) UpdateRPTypeService(requestParams *request.UpdateRPTypeRequest) (string, error) {
	db, err := database.RPConnection()
	if err != nil {
		return "error", fmt.Errorf("database connection error: %v", err)
	}
	dbInstance, _ := db.DB()

	query := fmt.Sprintf(`UPDATE report_mng SET rptid = '%d' WHERE rpcode = '%s'`, requestParams.Rptid, requestParams.Rpcode)

	err = db.Exec(query).Error
	if err != nil {
		return "error", err
	}

	dbInstance.Close()

	return "success", nil

}

// func - ExcelByDepartmantsService
func (s *AdminService) ExcelByDepartmentsService() ([]types.ExcelByDepartmantsService, error) {
	var excelByDepartments []types.ExcelByDepartmantsService
	db, err := database.RPConnection()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return nil, err
	}
	dbInstance, _ := db.DB()

	// Truy vấn lưu dữ liệu
	query := `
WITH MonthlyReports AS (
    SELECT
        rm.dpid,
        dep.title,
        MONTH(CAST(rm.updatedate AS DATE)) AS ReportMonth,
        COUNT(*) AS ReportCount,
        SUM(CASE WHEN rm.is_completed = 1 THEN 1 ELSE 0 END) AS CompletedCount
    FROM report_mng rm
    LEFT JOIN report_departments dep ON dep.dpid = rm.dpid
    WHERE YEAR(CAST(rm.updatedate AS DATE)) = YEAR(GETDATE())
    GROUP BY rm.dpid, dep.title, MONTH(CAST(rm.updatedate AS DATE))
),
PivotedData AS (
    SELECT *
    FROM MonthlyReports
    PIVOT (
        SUM(ReportCount)
        FOR ReportMonth IN ([1], [2], [3], [4], [5], [6], [7], [8], [9], [10], [11], [12])
    ) AS PivotTable
),
TotalReports AS (
    SELECT SUM(ReportCount) AS TotalAllReports
    FROM MonthlyReports
),
CompletedReports AS (
    SELECT 
        dpid,
        SUM(CompletedCount) AS TotalCompleted
    FROM MonthlyReports
    GROUP BY dpid
),
-- Thêm ROW_NUMBER() để chỉ lấy TotalCompleted cho dòng đầu tiên của mỗi dpid
DepartmentStats AS (
    SELECT 
        p.dpid,
        p.title,
        ISNULL([1], 0) AS Jan, ISNULL([2], 0) AS Feb, ISNULL([3], 0) AS Mar,
        ISNULL([4], 0) AS Apr, ISNULL([5], 0) AS May, ISNULL([6], 0) AS Jun,
        ISNULL([7], 0) AS Jul, ISNULL([8], 0) AS Aug, ISNULL([9], 0) AS Sep,
        ISNULL([10], 0) AS Oct, ISNULL([11], 0) AS Nov, ISNULL([12], 0) AS Dec,
        (ISNULL([1], 0) + ISNULL([2], 0) + ISNULL([3], 0) + 
         ISNULL([4], 0) + ISNULL([5], 0) + ISNULL([6], 0) + 
         ISNULL([7], 0) + ISNULL([8], 0) + ISNULL([9], 0) + 
         ISNULL([10], 0) + ISNULL([11], 0) + ISNULL([12], 0)) AS TotalReports,
        CAST((ISNULL([1], 0) + ISNULL([2], 0) + ISNULL([3], 0) + 
              ISNULL([4], 0) + ISNULL([5], 0) + ISNULL([6], 0) + 
              ISNULL([7], 0) + ISNULL([8], 0) + ISNULL([9], 0) + 
              ISNULL([10], 0) + ISNULL([11], 0) + ISNULL([12], 0)) * 100.0 / 
             (SELECT TotalAllReports FROM TotalReports) AS DECIMAL(5,2)) AS Percentage,
        -- Chỉ giữ giá trị TotalCompleted cho dòng đầu tiên (ROW_NUMBER() = 1)
        CASE 
            WHEN ROW_NUMBER() OVER (PARTITION BY p.dpid ORDER BY p.dpid) = 1 
            THEN cr.TotalCompleted ELSE 0 
        END AS TotalCompleted
    FROM PivotedData p
    JOIN CompletedReports cr ON p.dpid = cr.dpid
),
FinalResult AS (
    SELECT 
        title AS Department, 
        SUM(Jan) AS Jan, SUM(Feb) AS Feb, SUM(Mar) AS Mar, 
        SUM(Apr) AS Apr, SUM(May) AS May, SUM(Jun) AS Jun, 
        SUM(Jul) AS Jul, SUM(Aug) AS Aug, SUM(Sep) AS Sep, 
        SUM(Oct) AS Oct, SUM(Nov) AS Nov, SUM(Dec) AS Dec, 
        SUM(TotalReports) AS TotalReports, 
        SUM(Percentage) AS Percentage,
        CAST(CASE 
            WHEN SUM(TotalReports) = 0 THEN 0 
            ELSE (SUM(TotalCompleted) * 100.0 / SUM(TotalReports)) 
        END AS DECIMAL(5,2)) AS CompletionPercentage,
        MIN(CASE WHEN title = 'TC' THEN 1 ELSE 0 END) AS SortOrder,
        MIN(dpid) AS MinDpid
    FROM DepartmentStats
    GROUP BY title

    UNION ALL

    SELECT 
        'TC' AS Department,
        SUM(Jan) AS Jan, SUM(Feb) AS Feb, SUM(Mar) AS Mar,
        SUM(Apr) AS Apr, SUM(May) AS May, SUM(Jun) AS Jun,
        SUM(Jul) AS Jul, SUM(Aug) AS Aug, SUM(Sep) AS Sep,
        SUM(Oct) AS Oct, SUM(Nov) AS Nov, SUM(Dec) AS Dec,
        SUM(TotalReports) AS TotalReports,
        100.00 AS Percentage,
        CAST(SUM(TotalCompleted) * 100.0 / SUM(TotalReports) AS DECIMAL(5,2)) AS CompletionPercentage,
        1 AS SortOrder,
        999999 AS MinDpid  -- A large number to ensure TC is always last
    FROM DepartmentStats
)
SELECT 
    Department, Jan, Feb, Mar, Apr, May, Jun, Jul, Aug, Sep, Oct, Nov, Dec,
    TotalReports, Percentage, CompletionPercentage
FROM FinalResult
ORDER BY 
    SortOrder,
    MinDpid
	`
	err = db.Raw(query).Scan(&excelByDepartments).Error
	if err != nil {
		fmt.Println("Query error:", err)
		return nil, err
	}
	dbInstance.Close()
	return excelByDepartments, nil

}

// func - ExcelByReportTypesService
func (s *AdminService) ExcelByReportTypesService() ([]types.ExcelByReportTypesService, error) {
	var excelByReportTypes []types.ExcelByReportTypesService
	db, err := database.RPConnection()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return nil, err
	}
	dbInstance, _ := db.DB()
	// truy vấn lưu trững dữ liệu
	query := `
	WITH MonthlyReports AS (
    SELECT
        rm.rptid,
        rpt.title,
        MONTH(CAST(rm.updatedate AS DATE)) AS ReportMonth,
        COUNT(*) AS ReportCount,
        SUM(CASE WHEN rm.is_completed = 1 THEN 1 ELSE 0 END) AS CompletedCount
    FROM report_mng rm
    LEFT JOIN report_types rpt ON rpt.rptid = rm.rptid
    WHERE YEAR(CAST(rm.updatedate AS DATE)) = YEAR(GETDATE())
    GROUP BY rm.rptid, rpt.title, MONTH(CAST(rm.updatedate AS DATE))
),
PivotedData AS (
    SELECT *
    FROM MonthlyReports
    PIVOT (
        SUM(ReportCount)
        FOR ReportMonth IN ([1], [2], [3], [4], [5], [6], [7], [8], [9], [10], [11], [12])
    ) AS PivotTable
),
TotalReports AS (
    SELECT SUM(ReportCount) AS TotalAllReports
    FROM MonthlyReports
),
CompletedReports AS (
    SELECT 
        rptid,
        SUM(CompletedCount) AS TotalCompleted,
        SUM(ReportCount) AS TotalReports
    FROM MonthlyReports
    GROUP BY rptid
),
ReportTypeStats AS (
    SELECT 
        p.rptid,
        MAX(p.title) AS title, -- Use MAX to get a single title for each rptid
        SUM(ISNULL([1], 0)) AS Jan, SUM(ISNULL([2], 0)) AS Feb, SUM(ISNULL([3], 0)) AS Mar,
        SUM(ISNULL([4], 0)) AS Apr, SUM(ISNULL([5], 0)) AS May, SUM(ISNULL([6], 0)) AS Jun,
        SUM(ISNULL([7], 0)) AS Jul, SUM(ISNULL([8], 0)) AS Aug, SUM(ISNULL([9], 0)) AS Sep,
        SUM(ISNULL([10], 0)) AS Oct, SUM(ISNULL([11], 0)) AS Nov, SUM(ISNULL([12], 0)) AS Dec,
        MAX(cr.TotalReports) AS TotalReports,
        CAST(MAX(cr.TotalReports) * 100.0 / 
             (SELECT TotalAllReports FROM TotalReports) AS DECIMAL(5,2)) AS Percentage,
MAX(cr.TotalCompleted) AS TotalCompleted
    FROM PivotedData p
    JOIN CompletedReports cr ON p.rptid = cr.rptid
    GROUP BY p.rptid
),
FinalResult AS (
    SELECT 
        title AS ReportType, 
        SUM(Jan) AS Jan, SUM(Feb) AS Feb, SUM(Mar) AS Mar, 
        SUM(Apr) AS Apr, SUM(May) AS May, SUM(Jun) AS Jun, 
        SUM(Jul) AS Jul, SUM(Aug) AS Aug, SUM(Sep) AS Sep, 
        SUM(Oct) AS Oct, SUM(Nov) AS Nov, SUM(Dec) AS Dec, 
        SUM(TotalReports) AS TotalReports, 
        SUM(Percentage) AS Percentage,
        CAST(CASE 
            WHEN SUM(TotalReports) = 0 THEN 0 
            ELSE (SUM(TotalCompleted) * 100.0 / SUM(TotalReports)) 
        END AS DECIMAL(5,2)) AS CompletionPercentage,
        MIN(CASE WHEN title = 'TC' THEN 1 ELSE 0 END) AS SortOrder,
        MIN(rptid) AS MinRptid
    FROM ReportTypeStats
    GROUP BY title
    UNION ALL
    SELECT 
        'TC' AS ReportType,
        SUM(Jan) AS Jan, SUM(Feb) AS Feb, SUM(Mar) AS Mar,
        SUM(Apr) AS Apr, SUM(May) AS May, SUM(Jun) AS Jun,
        SUM(Jul) AS Jul, SUM(Aug) AS Aug, SUM(Sep) AS Sep,
        SUM(Oct) AS Oct, SUM(Nov) AS Nov, SUM(Dec) AS Dec,
        SUM(TotalReports) AS TotalReports,
        100.00 AS Percentage,
        CAST(SUM(TotalCompleted) * 100.0 / SUM(TotalReports) AS DECIMAL(5,2)) AS CompletionPercentage,
        1 AS SortOrder,
        999999 AS MinRptid
    FROM ReportTypeStats
)
SELECT 
    ReportType, Jan, Feb, Mar, Apr, May, Jun, Jul, Aug, Sep, Oct, Nov, Dec,
    TotalReports, Percentage, CompletionPercentage
FROM FinalResult
ORDER BY 
    SortOrder,
    MinRptid
	`
	err = db.Raw(query).Scan(&excelByReportTypes).Error
	if err != nil {
		fmt.Println("Query error:", err)
		return nil, err
	}
	dbInstance.Close()
	return excelByReportTypes, nil
}
