package output

type Output struct {
	ID               int
	EmployeeNumber   string
	FamilyName       string
	GivenName        string
	FamilyNameKana   string
	GivenNameKana    string
	Email            string
	EmploymentStatus string
	Department       RelatedOutput
	Position         RelatedOutput
}

type RelatedOutput struct {
	ID   int
	Code string
	Name string
}
