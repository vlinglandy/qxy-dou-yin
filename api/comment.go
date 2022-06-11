package api

import (
	"github.com/gin-gonic/gin"
	"qxy-dy/model"
	"qxy-dy/serializer"
	"strconv"
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
				CreateDate: commentList[i].CreateDate,
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
