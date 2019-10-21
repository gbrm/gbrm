package model


type Users struct {
	Id        int `gorm:"primary_key" json:"id"`
	Name  string              `gorm:"unique_index, not null" json:"Name"`
	Email  string             `gorm:"unique_index, not null" json:"Email"`
	Password string           `json:"Password"`
	Token string              `json:"Token"`
	IsAuth bool               `json:"IsAuth"`
	ProTill int               `json:"ProTill"`
	Type int                  `json:"Type"`
	Messages []Messages       `gorm:"ForeignKey:ToUserId;ASSOCIATION_FOREIGNKEY:Id" json:"messages"`
}
