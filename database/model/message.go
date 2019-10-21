package model

type Messages struct {
	Id        int `gorm:"primary_key" json:"id"`
	FromName    string `json:"from_name"`
	SentAt      int    `json:"sent_at"`
	Content     string `gorm:"type:longText" json:"content"`
	MessageType int    `json:"message_type"`
	IsRead      bool    `json:"is_read"`
	ToUserId    int    `json:"to_user_id"`
}
