package api

import (
	"fmt"
	"qxy-dy/model"
	"qxy-dy/serializer"

	"github.com/gin-gonic/gin"
)

// 这个文件是api的示例文件，大家编写代码可以参考

type UserResponse struct {
	serializer.Response
	ID            uint   `json:"id"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
}

func Demo(c *gin.Context) {
	var user model.User
	id := c.Query("id")
	if err := model.DB.First(&user, id).Error; err == nil {
		fmt.Printf("user:%#v", user)
		c.JSON(200, UserResponse{
			Response:      serializer.Response{StatusCode: 0},
			ID:            user.ID,
			Username:      user.Username,
			Password:      user.Password,
			FollowCount:   user.FollowCount,
			FollowerCount: user.FollowerCount,
		})
	} else {
		c.JSON(200, err)
	}
}
