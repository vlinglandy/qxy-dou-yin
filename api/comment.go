package api

import (
	"qxy-dy/serializer"

	"github.com/gin-gonic/gin"
)

// 评论操作
func ComentAction(c *gin.Context) {
	c.JSON(200, serializer.Response{
		Data: "", // 相应数据
		Code: 0,
		Msg:  "在这里写提示消息",
	})
}

// 视频评论列表
func CommentList(c *gin.Context) {
	c.JSON(200, serializer.Response{
		Data: "", // 相应数据
		Code: 0,
		Msg:  "在这里写提示消息",
	})
}
