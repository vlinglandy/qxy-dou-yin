package api

import (
	"fmt"
	"qxy-dy/model"
	"qxy-dy/serializer"

	"github.com/gin-gonic/gin"
)

// 这个文件是api的示例文件，大家编写代码可以参考

// 首先定义返回类型，如果返回中的实体已经被定义了就可以直接用写好的实体

type UserResponseDemo struct {
	serializer.Response
	ID            uint   `json:"id"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
}

func Demo(c *gin.Context) {
	var user model.User
	// 获取参数
	id := c.Query("id")
	// 操作数据库，可以在model对应的表中写好之后调用，也可以直接在这里操作数据库，不过直接操作不规范
	if err := model.DB.First(&user, id).Error; err == nil {
		fmt.Printf("user:%#v", user)
		// 利用已经定义好的响应结构体（一定要和接口文档的参数对应），返回响应数据
		c.JSON(200, UserResponseDemo{
			Response:      serializer.Response{StatusCode: 0},
			ID:            user.ID,
			Username:      user.Username,
			Password:      user.Password,
			FollowCount:   user.FollowCount,
			FollowerCount: user.FollowerCount,
		})
	} else {
		// 处理错误信息，基本上不用动，直接cv
		c.JSON(200, serializer.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	}
}
