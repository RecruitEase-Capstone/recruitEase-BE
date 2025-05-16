package model

type CVSummarizeResponse struct {
	Name              []string `json:"name"`
	CollegeName       []string `json:"college_name"`
	Degree            []string `json:"degree"`
	GraduationYear    []string `json:"graduation_year"`
	YearsOfExperience []string `json:"years_of_experience"`
	CompaniesWorkedAt []string `json:"companies_worked_at"`
	Designation       []string `json:"designation"`
	Skills            []string `json:"skills"`
	Location          []string `json:"location"`
	EmailAddress      []string `json:"email_address"`
}
