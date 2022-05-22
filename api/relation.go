package api

import (
	"qxy-dy/serializer"

	"github.com/gin-gonic/gin"
)

// 关系操作
func RelationAction(c *gin.Context) {
	c.JSON(200, serializer.Response{
		Data: "", // 相应数据
		Code: 0,
		Msg:  "在这里写提示消息",
	})
}

// 用户关注列表
func UserFollowList(c *gin.Context) {
	c.JSON(200, serializer.Response{
		Data: "", // 相应数据
		Code: 0,
		Msg:  "在这里写提示消息",
	})
}

// 评论操作
func UserLoverList(c *gin.Context) {
	c.JSON(200, serializer.Response{
		Data: "", // 相应数据
		Code: 0,
		Msg:  "在这里写提示消息",
	})
}
