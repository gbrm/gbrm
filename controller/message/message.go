package message

import (
	"github.com/shouc/gbrm/database"
	"github.com/shouc/gbrm/database/model"
	"github.com/shouc/gbrm/utils"
	"github.com/gin-gonic/gin"
)

func UnreadMessages(c *gin.Context) {
	// Auth Section
	isAuth := c.GetBool("isAuth")
	if !isAuth {
		return
	}

	// Message Section
	var messages []model.Messages
	cursor := database.GetDB()
	userId := utils.GetUserIdFromContext(c)
	var total int
	cursor.

		Table("messages").
		Where("to_user_id = ?", userId).
		Where("is_read = 0").
		Find(&messages)
	c.Set("total", total)
	c.Set("result", messages)
	c.Set("success", 1)
}

func ReadMessages(c *gin.Context) {
	// Auth Section
	isAuth := c.GetBool("isAuth")
	if !isAuth {
		return
	}

	// Message Section
	var messages []model.Messages
	cursor := database.GetDB()
	userId := utils.GetUserIdFromContext(c)
	offset, limit := utils.GetPageFromContext(c)
	var total int
	cursor.

		Table("messages").
		Where("to_user_id = ?", userId).
		Where("is_read = 1").
		Order("sent_at DESC").
		Count(&total).
		Offset(offset).
		Limit(limit).
		Find(&messages)
	c.Set("total", total)
	c.Set("result", messages)
	c.Set("success", 1)
}

func AddMessage(
	fromName string,
	toUserId int,
	content string,
	messageType int) {
	var newMessage = model.Messages{
		FromName:    fromName,
		ToUserId:    toUserId,
		Content:     content,
		MessageType: messageType,
		SentAt:      utils.CurrentTime(),
		IsRead:      false,
	}
	cursor := database.GetDB()
	cursor.Create(&newMessage)
}
