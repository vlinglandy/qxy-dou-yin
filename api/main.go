package api

import (
	"encoding/json"
	"qxy-dy/serializer"
	"strings"

	"github.com/gin-gonic/gin"
)

// Ping
func Ping(c *gin.Context) {
	c.JSON(200, serializer.Response{
		StatusCode: 0,
		StatusMsg:  "pong",
	})
}

// CurrentUser 获取当前用户id
func CurrentUser(c *gin.Context) string {
	value, _ := c.Get("MyId")
	encoded, _ := json.Marshal(value)
	myId := string(encoded)
	return strings.Trim(myId, "\"")
}

// ErrorResponse 返回错误消息
func ErrorResponse(err error) {

}
