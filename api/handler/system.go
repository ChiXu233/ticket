package handler

import (
	"github.com/gin-gonic/gin"
	config "ticket-service/conf"
)

func (handler *RestHandler) Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"code":    200,
		"message": config.Conf.DB.Name,
	})
}
