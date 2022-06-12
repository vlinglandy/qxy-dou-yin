package api

import (
	"net/http"
	"qxy-dy/middleware"
	"qxy-dy/model"
	"qxy-dy/serializer"
	"qxy-dy/util"

	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	serializer.Response
	VideoList []serializer.Video `json:"video_list"`
	NextTime  int64              `json:"next_time"`
}
type FeedService struct{}

func (FeedService) Feed(userId uint64, latestTime int64) []serializer.Video {
	var videos []model.Video
	var Videos []serializer.Video
	model.DB.Model(&model.Video{}).Where("created_at > 0").Order("created_at DESC").Limit(30).Find(&videos)
	ip := util.GetIp()
	for _, video := range videos {
		var author serializer.User
		// 通过AuthorId查找视频作者相关信息
		model.DB.Model(&serializer.User{}).Where("id = ?", video.AuthorId).First(&author)

		var isFollowCount, isFavoriteCount int64
		if userId > 0 { // 解析出了用户id
			// 查询是否已follow
			model.DB.Model(&model.Follow{}).Where("to_user_id = ? and user_id = ?", author.Id, userId).Count(&isFollowCount)
			author.IsFollow = isFollowCount > 0

			// 查询是否已喜欢
			model.DB.Model(&model.Favorite{}).Where("user_id = ? and video_id = ?", userId, video.ID).Count(&isFavoriteCount)
		}

		// 视频点赞数目
		var FavoriteCount int64
		model.DB.Model(&model.Favorite{}).Where("video_id = ?",video.ID).Count(&FavoriteCount)

		// 视频评论数目
		var CommentCount int64
		model.DB.Model(&model.Comment{}).Where("video_id = ?",video.ID).Count(&CommentCount)

		playurl, _ := util.ReplaceIP(video.PlayUrl, ip)
		coverurl, _ := util.ReplaceIP(video.CoverUrl, ip)
		Videos = append(Videos, serializer.Video{
			Id:            int64(video.ID),
			Author:        author,
			PlayUrl:       playurl,
			CoverUrl:      coverurl,
			FavoriteCount: FavoriteCount,
			CommentCount:  CommentCount,
			IsFavorite:    isFavoriteCount > 0,
			Title:         video.Title,
		})
	}
	return Videos
}

var feedService FeedService

// Feed 无需登录，返回按投稿时间倒序的视频列表，视频数由服务端控制，单次最多30个。
func GetVideo(c *gin.Context) {
	var latestTimeTemp = c.Query("latest_time")
	latestTime, _ := strconv.ParseInt(latestTimeTemp, 10, 64)
	// 获取token
	token := c.Query("token")

	// 通过token获取id, 生成失败时id为0，数据库中无id为0的数据
	id, _ := middleware.GetIdByToken(token)

	// 如果返回长度大于0则表示有视频，返回给用户
	if Videos := feedService.Feed(id, latestTime); len(Videos) > 0 {
		c.JSON(http.StatusOK, serializer.FeedResponse{
			Response: serializer.Response{
				StatusCode: 0,
				StatusMsg:  "获取视频流成功",
			},
			VideoList: Videos,
			NextTime:  time.Now().Unix(),
		})
		return
	}

	// 其他情况，返回UnknownError
	c.JSON(http.StatusOK, serializer.Response{
		StatusCode: 1,
		StatusMsg:  "未知错误",
	})
}
