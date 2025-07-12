package models

// Profile represents a user's profile information
type Profile struct {
	BaseModelWithUser

	Address    string `json:"address"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	PostalCode string `json:"postal_code"`

	CompanyName string `json:"company_name"`
	JobTitle    string `json:"job_title"`
	Department  string `json:"department"`
}
