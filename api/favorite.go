package api

import (
	"qxy-dy/serializer"

	"github.com/gin-gonic/gin"
)

// 点赞操作
func FavoriteAction(c *gin.Context) {
	c.JSON(200, serializer.Response{
		Data: "", // 相应数据
		Code: 0,
		Msg:  "在这里写提示消息",
	})
}

// 点赞列表
func FavoriteList(c *gin.Context) {
	c.JSON(200, serializer.Response{
		Data: "", // 相应数据
		Code: 0,
		Msg:  "在这里写提示消息",
	})
}
