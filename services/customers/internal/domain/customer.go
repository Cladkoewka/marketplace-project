package domain

type Customer struct {
	ID            int64  `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"json"`
	Phone         string `json:"phone"`
	Address       string `json:"address"`
	ActivityScore int64    `json:"activity_score"`
}
