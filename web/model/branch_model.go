package model

type CreateManagementRequest struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Email   string `json:"email"`
}

type CreateManagementResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

type CreateBranchRequest struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Email   string `json:"email"`
	Upline  string `json:"upline"`
}

type CreateBranchResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}
