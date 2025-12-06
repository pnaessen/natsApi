package models

type User42 struct {
	Id        int    `json:"id"`
	Username    string `json:"login"`
	Email       string `json:"email"`
	School_year string `json:"pool_year"`
	Is_active      bool   `json:"active?"`
}
