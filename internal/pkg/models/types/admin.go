package types

// FILE - TYPES ADMIN

// type - Admin
type Admin struct {
	Adminid  uint   `json:"adminid" form:"adminid"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

// type - ExcelByDepartmants
type ExcelByDepartmantsService struct {
	Department           string `json:"department" form:"department"`
	Jan                  string `json:"Jan" form:"Jan"`
	Feb                  string `json:"Feb" form:"Feb"`
	Mar                  string `json:"Mar" form:"Mar"`
	Apr                  string `json:"Apr" form:"Apr"`
	May                  string `json:"May" form:"May"`
	Jun                  string `json:"Jun" form:"Jun"`
	Jul                  string `json:"Jul" form:"Jul"`
	Aug                  string `json:"Aug" form:"Aug"`
	Sep                  string `json:"Sep" form:"Sep"`
	Oct                  string `json:"Oct" form:"Oct"`
	Nov                  string `json:"Nov" form:"Nov"`
	Dec                  string `json:"Dec" form:"Dec"`
	TotalReports         string `json:"TotalReports" form:"TotalReports"`
	Percentage           string `json:"Percentage" form:"Percentage"`
	CompletionPercentage string `json:"CompletionPercentage" form:"CompletionPercentage"`
}

// type - ExcelByRpType
type ExcelByReportTypesService struct {
	ReportType           string `json:"ReportType" form:"ReportType"`
	Jan                  string `json:"Jan" form:"Jan"`
	Feb                  string `json:"Feb" form:"Feb"`
	Mar                  string `json:"Mar" form:"Mar"`
	Apr                  string `json:"Apr" form:"Apr"`
	May                  string `json:"May" form:"May"`
	Jun                  string `json:"Jun" form:"Jun"`
	Jul                  string `json:"Jul" form:"Jul"`
	Aug                  string `json:"Aug" form:"Aug"`
	Sep                  string `json:"Sep" form:"Sep"`
	Oct                  string `json:"Oct" form:"Oct"`
	Nov                  string `json:"Nov" form:"Nov"`
	Dec                  string `json:"Dec" form:"Dec"`
	TotalReports         string `json:"TotalReports" form:"TotalReports"`
	Percentage           string `json:"Percentage" form:"Percentage"`
	CompletionPercentage string `json:"CompletionPercentage" form:"CompletionPercentage"`
}
