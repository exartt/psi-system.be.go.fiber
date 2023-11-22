package person

type PersonPatient struct {
	ID       uint   `gorm:"column:id"`
	Name     string `gorm:"column:name"`
	Email    string `gorm:"column:email"`
	IsActive bool   `gorm:"column:is_active"`
	IsPlan   bool   `gorm:"column:is_plan"`
}

type DTO struct {
	ID             uint   `json:"ID"`
	PersonId       uint   `json:"PersonId"`
	Name           string `json:"Name"`
	Email          string `json:"Email"`
	CellPhone      string `json:"CellPhone"`
	Phone          string `json:"Phone"`
	ZipCode        string `json:"ZipCode"`
	Address        string `json:"Address"`
	IsActive       bool   `json:"IsActive"`
	CPF            string `json:"CPF"`
	RG             string `json:"RG"`
	IsPlan         bool   `json:"IsPlan"`
	SessionPrice   string `json:"SessionPrice"`
	ConversionType string `json:"ConversionType"`
}
