package api

import (
	"fmt"
	"qxy-dy/model"
	"qxy-dy/serializer"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FollowResponse struct {
	serializer.Response
	UserList []serializer.User `json:"user_list"`
}

// 关系操作
func RelationAction(c *gin.Context) {
	myId, _ := strconv.ParseUint(CurrentUser(c), 10, 32)
	id, _ := strconv.ParseUint(c.Query("to_user_id"), 10, 32)
	actionType := c.Query("action_type")
	fmt.Printf("myid:%#v", myId)
	follow := model.Follow{
		UserId:   uint(myId),
		ToUserId: uint(id),
	}
	var followCount int64 = 0
	model.DB.Model(&model.Follow{}).Where("user_id = ? AND to_user_id = ?", myId, id).Count(&followCount)
	if actionType == "1" {
		if followCount == 0 {
			result := model.DB.Model(&model.Follow{}).Create(&follow)
			if result.RowsAffected > 0 {
				c.JSON(200, serializer.Response{
					StatusCode: 0,
					StatusMsg:  "关注成功",
				})
			} else {
				c.JSON(200, serializer.Response{
					StatusCode: 1,
					StatusMsg:  "关注失败",
				})
			}
		} else {
			c.JSON(200, serializer.Response{
				StatusCode: 1,
				StatusMsg:  "您已经关注过了",
			})
		}

	} else if actionType == "2" {
		if followCount > 0 {
			deleteRes := model.DB.Model(&model.Follow{}).Where("user_id = ? AND to_user_id = ?", myId, id).Delete(&follow)
			if deleteRes.RowsAffected > 0 {
				c.JSON(200, serializer.Response{
					StatusCode: 0,
					StatusMsg:  "取关成功",
				})
			} else {
				c.JSON(200, serializer.Response{
					StatusCode: 1,
					StatusMsg:  "取关失败",
				})
			}
		} else {
			c.JSON(200, serializer.Response{
				StatusCode: 0,
				StatusMsg:  "您未关注该用户",
			})
		}

	} else {
		c.JSON(200, serializer.Response{
			StatusCode: 1,
			StatusMsg:  "类型无效",
		})
	}

}

// 用户关注列表MyUserService
func UserFollowList(c *gin.Context) {
	// id := c.Query("user_id")
	// model.DB.Model(&model.Follow{}).Where("user_id = ?", id).Find()

}

// 评论操作
func UserLoverList(c *gin.Context) {

}
