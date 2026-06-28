package input

// InputはCLIから受け取った社員検索条件をServiceへ渡すための型。
type Input struct {
	Keyword          string
	DepartmentID     string
	PositionID       string
	EmploymentStatus string
}
