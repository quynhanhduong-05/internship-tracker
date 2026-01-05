package application

type CreateRequest struct {
	CompanyName string `json:"company_name"`
	Position    string `json:"position"`
	Status      string `json:"status"`
	Notes       string `json:"notes"`
}

type UpdateStatusRequest struct {
	Status string `json:"status"`
}
