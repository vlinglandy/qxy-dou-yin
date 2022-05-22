package api

import (
	"qxy-dy/serializer"

	"github.com/gin-gonic/gin"
)

// 视频投稿
func PublishAction(c *gin.Context) {
	c.JSON(200, serializer.Response{
		Data: "", // 相应数据
		Code: 0,
		Msg:  "在这里写提示消息",
	})
}

// 自己发布视频列表
func PubishList(c *gin.Context) {
	c.JSON(200, serializer.Response{
		Data: "", // 相应数据
		Code: 0,
		Msg:  "在这里写提示消息",
	})
}
