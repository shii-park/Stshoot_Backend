package model

type Message struct {
	Username string `json:"username"`
	Text     string `json:"text"`
	Price    int    `json:"price"`
}
