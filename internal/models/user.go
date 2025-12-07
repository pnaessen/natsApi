package models

type User42 struct {
	ID          int    `json:"id"`
	Username    string `json:"login"`
	Email       string `json:"email"`
	School_year string `json:"pool_year"`
	Is_active   bool   `json:"active?"`
}

type UserMessage struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	IntraID    int    `json:"intra_id"`
	SchoolYear string `json:"school_year"`
	IsActive   bool   `json:"is_active"`

	Db_id     uint   `json:"db_id,omitempty"`
}
