package model

type Todos struct {
	Id          int    `json:"id"`
	UserId      string `json:"user_id"`
	Description string `json:"description"`
	Title       string `json:"title"`
}