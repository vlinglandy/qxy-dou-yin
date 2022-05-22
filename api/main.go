package api

import (
	"qxy-dy/serializer"

	"github.com/gin-gonic/gin"
)

// Ping
func Ping(c *gin.Context) {
	c.JSON(200, serializer.Response{
		StatusCode: 0,
		StatusMsg:  "pong",
	})
}

// CurrentUser 获取当前用户
func CurrentUser(c *gin.Context) {

}

// ErrorResponse 返回错误消息
func ErrorResponse(err error) {

}
