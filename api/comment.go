package api

import (
	"github.com/gin-gonic/gin"
	"qxy-dy/model"
	"qxy-dy/serializer"
	"strconv"
	"time"
)

//type Comment struct {
//	gorm.Model
//	UserId     uint
//	Content    string
//	CreateDate string
//	VideoId    uint
//}

type CommentResponse struct {
	serializer.Response
	CommentList []serializer.Comment `json:"comment_list"`
}

// 评论操作
func ComentAction(c *gin.Context) {
	uid, _ := strconv.ParseUint(CurrentUser(c), 10, 32)
	video_id, _ := strconv.ParseUint(c.Query("video_id"), 10, 32)
	var commentResponse []serializer.Comment
	actionType := c.Query("action_type")
	if actionType == "1" {
		commentText := c.Query("comment_text")
		comment := model.Comment{
			UserId:     uint(uid),
			Content:    commentText,
			CreateDate: time.Now().Format("2006-01-02 15:04:05"),
			VideoId:    uint(video_id),
		}
		result := model.DB.Model(&model.Comment{}).Create(&comment)
		if result.RowsAffected > 0 {
			isFollowed, _ := MyUserService.isFollowByUint(CurrentUser(c), uint(uid))
			user := model.User{}
			model.DB.Model(&model.User{}).Where("id = ?", uint(uid)).Find(&user)
			t, _ := time.ParseInLocation("2006-01-02 15:04:05", comment.CreateDate, time.Local)
			commentResponse = append(commentResponse, serializer.Comment{
				Id: int64(comment.ID),
				User: serializer.User{
					Id:            int64(user.ID),
					Name:          user.Username,
					FollowCount:   user.FollowCount,
					FollowerCount: user.FollowerCount,
					IsFollow:      isFollowed,
				},
				Content:    comment.Content,
				CreateDate: strconv.Itoa(int(t.Month())) + "-" + strconv.Itoa(t.Day()),
			})

			c.JSON(200, CommentResponse{
				Response: serializer.Response{
					StatusCode: 0,
					StatusMsg:  "评论成功",
				},
				CommentList: commentResponse,
			})
		} else {
			c.JSON(200, serializer.Response{
				StatusCode: 1,
				StatusMsg:  "评论失败",
			})
		}
	} else if actionType == "2" {
		commentId := c.Query("comment_id")
		comment := model.Comment{}
		id, _ := strconv.ParseUint(commentId, 10, 32)
		comment.ID = uint(id)
		result := model.DB.Model(&model.Comment{}).Where("id = ?", commentId).Delete(&comment)
		if result.RowsAffected > 0 {
			c.JSON(200, serializer.Response{
				StatusCode: 0,
				StatusMsg:  "删除评论成功",
			})
		} else {
			c.JSON(200, serializer.Response{
				StatusCode: 1,
				StatusMsg:  "删除评论失败",
			})
		}
	}
}

// 视频评论列表
func CommentList(c *gin.Context) {
	video_id, _ := strconv.ParseUint(c.Query("video_id"), 10, 32)
	var commentList []model.Comment
	var commentResponse []serializer.Comment
	model.DB.Model(&model.Comment{}).Where("video_id = ?", video_id).Find(&commentList)
	if len(commentList) > 0 {
		for i := 0; i < len(commentList); i++ {
			isFollowed, _ := MyUserService.isFollowByUint(CurrentUser(c), commentList[i].UserId)
			user := model.User{}
			model.DB.Model(&model.User{}).Where("id = ?", commentList[i].UserId).Find(&user)
			t, _ := time.ParseInLocation("2006-01-02 15:04:05", commentList[i].CreateDate, time.Local)
			commentResponse = append(commentResponse, serializer.Comment{
				Id: int64(commentList[i].ID),
				User: serializer.User{
					Id:            int64(user.ID),
					Name:          user.Username,
					FollowCount:   user.FollowCount,
					FollowerCount: user.FollowerCount,
					IsFollow:      isFollowed,
				},
				Content:    commentList[i].Content,
				CreateDate: strconv.Itoa(int(t.Month())) + "-" + strconv.Itoa(t.Day()),
			})
		}
		c.JSON(200, CommentResponse{
			Response: serializer.Response{
				StatusCode: 0,
				StatusMsg:  "成功获取评论列表",
			},
			CommentList: commentResponse,
		})
	} else {
		c.JSON(200, serializer.Response{
			StatusCode: 0,
			StatusMsg:  "没有评论",
		})
	}
}
