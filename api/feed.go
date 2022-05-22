package api

import (
	"qxy-dy/serializer"

	"github.com/gin-gonic/gin"
)

// 这个文件是api的示例文件，大家编写代码可以参考

func GetVideo(c *gin.Context) {
	c.JSON(200, serializer.Response{
		Data: "", // 相应数据
		Code: 0,
		Msg:  "在这里写提示消息",
	})
}
