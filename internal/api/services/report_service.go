package services

import (
	"database/sql"
	"fmt"
	"time"

	// "time"
	"web-api/internal/pkg/database"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/types"
)

// FILE - REPORT SERVICE
type ReportService struct {
	*BaseService
}

var RP = &ReportService{}

// func GetReportService
func (s *ReportService) GetReportService() ([]types.Report, error) {
	var reports []types.Report

	// kết nối cơ sở dữ liệu
	db, err := database.RPConnection()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return nil, err
	}
	dbInstance, _ := db.DB()
	// truy vấn lấy dữ liệu
	query := `SELECT * FROM reports`

	err = db.Raw(query).Scan(&reports).Error
	if err != nil {
		fmt.Println("Query error:", err)
		return nil, err
	}
	dbInstance.Close()
	return reports, nil
}

// func GetReportByRpcodeService
func (s *ReportService) GetReportByRpcodeService(requestParams *request.GetReportByRpcodeRequest) ([]types.Report, error) {
	var reports []types.Report
	var err error
	// kết nối cơ sở dữ liệu
	db, err := database.RPConnection()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return nil, err
	}
	dbInstance, _ := db.DB()
	// truy vấn lấy dữ liệu
	query := fmt.Sprintf(`SELECT * FROM reports WHERE rpcode = '%s'`, requestParams.Rpcode)
	err = db.Raw(query).Scan(&reports).Error
	if err != nil {
		fmt.Println("Query error:", err)
		return nil, err
	}
	dbInstance.Close()
	return reports, nil
}

// func InsertReportFirstService
func (s *ReportService) InsertReportFirstService(requestParams *request.InsertReportRequest) (string, error) {
	db, err := database.RPConnection()
	if err != nil {
		return "error", fmt.Errorf("database connection error: %v", err)
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Start a transaction
	tx := db.Begin()
	if tx.Error != nil {
		return "error", fmt.Errorf("failed to begin transaction: %v", tx.Error)
	}

	// Định dạng ngày giờ phù hợp với smalldatetime của SQL Server
	noteTime := time.Now().Format("2006-01-02 15:04:00")

	// Insert into reports table
	reportQuery := `INSERT INTO reports (note, rpcode, note_time) VALUES (@p1, @p2, @p3)`
	err = tx.Exec(reportQuery,
		sql.Named("p1", requestParams.Note),
		sql.Named("p2", requestParams.Rpcode),
		sql.Named("p3", noteTime)).Error
	if err != nil {
		tx.Rollback()
		return "error", fmt.Errorf("failed to insert into reports: %v", err)
	}

	// Insert into report_mng table
	reportMngQuery := `INSERT INTO report_mng (rpcode, dpid, updatedate) 
						VALUES (@p1, @p2, @p3)`
	err = tx.Exec(reportMngQuery,
		sql.Named("p1", requestParams.Rpcode),
		sql.Named("p2", requestParams.Dpid),
		sql.Named("p3", noteTime)).Error
	if err != nil {
		tx.Rollback()
		return "error", fmt.Errorf("failed to insert into report_mng: %v", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return "error", fmt.Errorf("failed to commit transaction: %v", err)
	}

	return "success", nil
}

func (s *ReportService) CheckRPCodeAndGetReportsService(requestParams *request.GetReportByRpcodeRequest) ([]types.CombinedReport, error) {
	db, err := database.RPConnection()
	if err != nil {
		return nil, fmt.Errorf("database connection error: %v", err)
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Check if rpcode exists in report_mng
	checkQuery := "SELECT COUNT(*) FROM report_mng WHERE rpcode = ?"
	var count int
	err = db.Raw(checkQuery, requestParams.Rpcode).Scan(&count).Error
	if err != nil {
		return nil, fmt.Errorf("error checking rpcode existence: %v", err)
	}

	if count == 0 {
		return nil, fmt.Errorf("rpcode %s does not exist", requestParams.Rpcode)
	}

	// If rpcode exists, fetch combined data from report_mng and reports
	combinedQuery := `
		SELECT 
			r.id, rm.rpcode, r.note, r.answer, rm.dpid, rm.rptid, 
			r.note_time, r.answer_time, rm.is_answered, rm.is_completed, rm.updatedate
		FROM 
			report_mng rm
		LEFT JOIN 
			reports r ON rm.rpcode = r.rpcode
		WHERE 
			rm.rpcode = ?
	`

	var combinedReports []types.CombinedReport
	err = db.Raw(combinedQuery, requestParams.Rpcode).Scan(&combinedReports).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching combined reports: %v", err)
	}

	return combinedReports, nil
}

// func - SendReportService
func (s *ReportService) SendReportService(requestParams *request.InsertReportRequest) (string, error) {
	db, err := database.RPConnection()
	if err != nil {
		return "error", fmt.Errorf("database connection error: %v", err)
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Start a transaction
	tx := db.Begin()
	if tx.Error != nil {
		return "error", fmt.Errorf("failed to begin transaction: %v", tx.Error)
	}

	noteTime := time.Now().Format("2006-01-02 15:04:05")

	// Insert into reports table
	insertQuery := `INSERT INTO reports (note, rpcode, note_time) VALUES (?, ?, ?)`
	err = tx.Exec(insertQuery, requestParams.Note, requestParams.Rpcode, noteTime).Error
	if err != nil {
		tx.Rollback()
		return "error", fmt.Errorf("failed to insert into reports: %v", err)
	}

	// Update report_mng table
	updateQuery := `UPDATE report_mng SET updatedate = ?, is_answered = 0 WHERE rpcode = ?`
	err = tx.Exec(updateQuery, noteTime, requestParams.Rpcode).Error
	if err != nil {
		tx.Rollback()
		return "error", fmt.Errorf("failed to update report_mng: %v", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return "error", fmt.Errorf("failed to commit transaction: %v", err)
	}

	return "success", nil
}

// func - UpdateNoteService
func (s *ReportService) UpdateNoteService(requestParams *request.UpdateNoteRequest) (string, error) {
	db, err := database.RPConnection()
	if err != nil {
		return "error", fmt.Errorf("database connection error: %v", err)
	}
	dbInstance, _ := db.DB()

	query := fmt.Sprintf(`UPDATE reports SET note = '%s' WHERE id = '%d' AND rpcode = '%s'`, requestParams.Note, requestParams.Id, requestParams.Rpcode)

	err = db.Exec(query).Error
	if err != nil {
		return "error", err
	}

	dbInstance.Close()

	return "success", nil
}

// func - UpdateIsCompletedService
func (s *ReportService) UpdateIsCompletedService(requestParams *request.GetReportByRpcodeRequest) (string, error) {
	db, err := database.RPConnection()
	if err != nil {
		return "error", fmt.Errorf("database connection error: %v", err)
	}
	dbInstance, _ := db.DB()

	query := fmt.Sprintf("UPDATE report_mng SET is_completed = 1, is_answered = 1 WHERE rpcode = '%s'", requestParams.Rpcode)

	err = db.Exec(query).Error
	if err != nil {
		return "error", err
	}

	dbInstance.Close()

	return "success", nil
}

// func - GetDepartmentsService
func (s *ReportService) GetDepartmentsService() ([]types.Department, error) {
	var departments []types.Department
	// Kết nối cơ sở dữ liệu
	db, err := database.RPConnection()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return nil, err
	}
	dbInstance, _ := db.DB()
	// Truy vấn lấy dữ liệu
	query := `SELECT * FROM report_departments`
	err = db.Raw(query).Scan(&departments).Error
	if err != nil {
		fmt.Println("Query error:", err)
		return nil, err
	}
	dbInstance.Close()
	return departments, nil
}

// func - SendEmailService
func (s *ReportService) SendEmailService() (string, error) {
	var emails []string

	// kết nối cơ sở dữ liệu
	db, err := database.RPConnection()
	if err != nil {
		return "error", fmt.Errorf("database connection error: %v", err)
	}
	dbInstance, _ := db.DB()

	// truy vấn lấy email
	query1 := "SELECT email FROM report_emails"
	err = db.Raw(query1).Scan(&emails).Error
	if err != nil {
		return "error", err
	}

	// gửi mail
	for _, email := range emails {
		query2 := "EXEC NotifyNewFeedbackEmail ?"
		result := db.Exec(query2, email)
		if result.Error != nil {
			return "error", result.Error
		}
	}

	dbInstance.Close()

	return "success", nil
}
