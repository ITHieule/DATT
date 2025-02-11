package request

//FILE - MODEL REQUEST REPORT
// request - InsertReportRequest
type InsertReportRequest struct {
	Rpcode string `json:"rpcode"`
	Dpid   uint   `json:"dpid"`
	Note   string `json:"note"`
}

// request - GetReportByRpcodeRequest
type GetReportByRpcodeRequest struct {
	Rpcode string `json:"rpcode"`
}


// request - AnswerReportRequest
type AnswerReportRequest struct {
	RPCode string `json:"rpcode"`
	Answer string `json:"answer"`
}

// request - EmailRequest
type EmailRequest struct {
	EmailType string `json:"email_type"`
}

// request - UpdateNoteRequest
type UpdateNoteRequest struct {
	Id     uint   `json:"id"`
	Note   string `json:"note"`
	Rpcode string `json:"rpcode"`
}
