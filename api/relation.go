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

// 用户关注列表
func UserFollowList(c *gin.Context) {
	var userFollowers []model.Follow
	var userResponse []serializer.User
	id := c.Query("user_id")
	model.DB.Model(&model.Follow{}).Where("user_id = ?", id).Find(&userFollowers)
	// fmt.Printf("userList:%#v\n", user)
	if len(userFollowers) > 0 {
		for i := 0; i < len(userFollowers); i++ {
			isFollowed, _ := MyUserService.isFollowByUint(CurrentUser(c), userFollowers[i].ToUserId)
			user := model.User{}
			model.DB.Model(&model.User{}).Where("id = ?", userFollowers[i].ToUserId).Find(&user)
			userResponse = append(userResponse, serializer.User{
				IsFollow:      isFollowed,
				Id:            int64(user.ID),
				Name:          user.Username,
				FollowCount:   user.FollowCount,
				FollowerCount: user.FollowerCount,
			})
		}
		c.JSON(200, FollowResponse{
			Response: serializer.Response{
				StatusCode: 0,
				StatusMsg:  "获取关注列表成功",
			},
			UserList: userResponse,
		})
	} else {
		c.JSON(200, serializer.Response{
			StatusCode: 0,
			StatusMsg:  "没有关注用户",
		})
	}
}

// 粉丝列表
func UserLoverList(c *gin.Context) {
	var userFollowers []model.Follow
	var userResponse []serializer.User
	id := c.Query("user_id")
	model.DB.Model(&model.Follow{}).Where("to_user_id = ?", id).Find(&userFollowers)
	// fmt.Printf("userList:%#v\n", user)
	if len(userFollowers) > 0 {
		for i := 0; i < len(userFollowers); i++ {
			isFollowed, _ := MyUserService.isFollowByUint(CurrentUser(c), userFollowers[i].UserId)
			user := model.User{}
			model.DB.Model(&model.User{}).Where("id = ?", userFollowers[i].UserId).Find(&user)
			userResponse = append(userResponse, serializer.User{
				IsFollow:      isFollowed,
				Id:            int64(user.ID),
				Name:          user.Username,
				FollowCount:   user.FollowCount,
				FollowerCount: user.FollowerCount,
			})
		}
		c.JSON(200, FollowResponse{
			Response: serializer.Response{
				StatusCode: 0,
				StatusMsg:  "获取粉丝列表成功",
			},
			UserList: userResponse,
		})
	} else {
		c.JSON(200, serializer.Response{
			StatusCode: 0,
			StatusMsg:  "没有粉丝",
		})
	}
}
