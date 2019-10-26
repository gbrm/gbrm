package model

type Petitions struct {
	Id          int
	InitiatorId int
	ReceiverId  int
	Content     string
	Status      int
}
