package modelgeneral

type LoginRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Pagination struct {
	Page        int `json:"page"`
	DataPerPage int `json:"dataPerPage"`
}
