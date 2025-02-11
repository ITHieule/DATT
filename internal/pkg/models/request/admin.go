package request

// FILE - REQUEST ADMIN

// request - LoginRequest
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// request - GetIsAnsweredbyReportTypeRequest
type GetReportByIsAnsweredAndRpidRequest struct {
	IsAnswered  bool `json:"is_answered"`
	Rptid       uint `json:"rptid"`
	IsCompleted bool `json:"is_completed"` // is_completed
}

// request - SendAnswerRequest
type SendAnswerRequest struct {
	Answer string `json:"answer"`
	Rpcode string `json:"rpcode"`
	Id     uint   `json:"id"`
}

// request - UpdateAnswerRequest
type UpdateAnswerRequest struct {
	Answer string `json:"answer"`
	Rpcode string `json:"rpcode"`
	Id     uint   `json:"id"`
}

// request - GetRPTypeRequest
type GetRPTTypeRequest struct {
	RPTID uint `json:"rptid"`
}

// request - UpdateRPTypeRequest
type UpdateRPTypeRequest struct {
	Rptid  uint   `json:"rptid"`
	Rpcode string `json:"rpcode"`
}
