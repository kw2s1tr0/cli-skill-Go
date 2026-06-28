package response

type Response struct {
	ID               int             `json:"id"`
	EmployeeNumber   string          `json:"employee_number"`
	DepartmentID     int             `json:"department_id"`
	PositionID       int             `json:"position_id"`
	FamilyName       string          `json:"family_name"`
	GivenName        string          `json:"given_name"`
	FamilyNameKana   string          `json:"family_name_kana"`
	GivenNameKana    string          `json:"given_name_kana"`
	Email            string          `json:"email"`
	EmploymentStatus string          `json:"employment_status"`
	CreatedAt        string          `json:"created_at"`
	UpdatedAt        string          `json:"updated_at"`
	Department       RelatedResponse `json:"department"`
	Position         RelatedResponse `json:"position"`
}

type RelatedResponse struct {
	ID        int    `json:"id"`
	Code      string `json:"code"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
