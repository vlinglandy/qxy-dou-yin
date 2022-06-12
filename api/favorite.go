package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"qxy-dy/middleware"
	"qxy-dy/model"
	"qxy-dy/serializer"
	"strconv"
	"time"
)

type LikesRsp struct {
	StatusCode 	int32	`json:"status_code"`
	StatusMsg	string	`json:"status_msg"`
}

type FavoriteListRsp struct{
	StatusCode 	int32				`json:"status_code"`
	StatusMsg	string				`json:"status_msg"`
	VideoList	[]serializer.Video	`json:"video_list"`
}


// FavoriteAction 点赞操作
func FavoriteAction(c *gin.Context) {
	// 登陆验证
	token := c.Query("token")
	videoId := c.Query("video_id")
	actionType := c.Query("action_type")

	// 状态码
	StatusCode := 0
	// 通过token获取id
	userId, err := middleware.GetIdByToken(token)
	if err != nil {
		StatusCode = 1
		c.JSON(http.StatusOK, serializer.Response{
			StatusCode: 1,
			StatusMsg:  "由token获取id失败",
		})
		return
	}

	// 转换一下格式
	videosId, _ := strconv.ParseInt(videoId,10, 32)

	// 构造点赞信息的实例
	favorites := model.Favorite{
		Model:   gorm.Model{},
		UserId: uint(userId),
		VideoId: uint(videosId),
	}

	// 用于接收查询的点赞信息
	var likesInfo model.Favorite

	// 查找数据库里有没有对应操作人 和点赞视频的消息
	if err := model.DB.Where("user_id = ? and video_id = ?", favorites.UserId, favorites.VideoId).Find(&likesInfo).Error; err != nil{
		StatusCode = 1
		c.JSON(200, LikesRsp{
			StatusCode: -1,
			StatusMsg: "点赞操作失败",
		})
	}

	//  actionType == "1" 点赞
	if actionType == "1" {
		// if: 刚刚的搜索没有结果
		if strconv.Itoa(int(likesInfo.UserId)) != strconv.FormatUint(userId, 10) {

			// 创建点赞信息
			if err =  model.DB.Create(&favorites).Error; err != nil{
				StatusCode = 1
				c.JSON(200, LikesRsp{
					StatusCode: -1,
					StatusMsg: "点赞操作失败",
				})
			}

		}else{
			StatusCode = 1
			// if: 刚刚的搜索有结果
			c.JSON(200, LikesRsp{
				StatusCode: -1,
				StatusMsg: "不可重复点赞",
			})

		}

	}else {
		// 如果是取消点赞
		if strconv.Itoa(int(likesInfo.UserId)) == strconv.FormatUint(userId, 10) {
			if err := model.DB.Model(likesInfo).Updates(map[string]interface{}{"deleted_at": time.Now(), "updated_at": time.Now()}).Error; err != nil {
				StatusCode = 1
				c.JSON(200, LikesRsp{
					StatusCode: -1,
					StatusMsg:  "取消点赞操作失败",
				})
			}
		}else{
			StatusCode = 1
			c.JSON(200, LikesRsp{
				StatusCode: -1,
				StatusMsg:  "不允许重复取消点赞",
			})
		}
	}

	if StatusCode == 0{
		c.JSON(200, LikesRsp{
			StatusCode: 0,
			StatusMsg: "操作成功",
		})
	}

}

// FavoriteList 点赞列表  视频的点赞列表
func FavoriteList(c *gin.Context) {

	// 用于记录状态
	// 传入的用户的id
	token := c.Query("token")
	userId := c.Query("user_id")

	// 状态码返回值
	StatusCode := 0

	// 通过token获取id
	user_id, err := middleware.GetIdByToken(token)
	if err != nil {
		StatusCode = 1
		c.JSON(http.StatusOK, serializer.Response{
			StatusCode: -1,
			StatusMsg:  "由token获取id失败",
		})
		return
	}

	// 校验token和传入id是不是同一用户
	if strconv.FormatUint(user_id, 10) != userId{
		StatusCode = 1
		c.JSON(200, FavoriteListRsp{
			StatusCode: -1,
			StatusMsg: "登录校验失败",
			VideoList: nil,
		})
	}

	// likeInfo: 存放点赞操作的信息
	var likeInfo []model.Favorite
	if err := model.DB.Where("user_id = ?" , userId).Find(&likeInfo).Error; err != nil{
		StatusCode = 1
		c.JSON(200, FavoriteListRsp{
			StatusCode: -1,
			StatusMsg: "操作失败",
			VideoList: nil,
		})
	}

	if len(likeInfo) == 0{
		StatusCode = 1
		c.JSON(200, FavoriteListRsp{
			StatusCode: 0,
			StatusMsg: "暂无点赞信息",
			VideoList: nil,
		})
	}

	// 通过视频id获得视频
	var videoList []model.Video
	for i := 0; i < len(likeInfo); i++{
		var video model.Video
		if err = model.DB.Where("id = ?" , likeInfo[i].VideoId).Find(&video).Error; err != nil{
			StatusCode = 1
			c.JSON(200, FavoriteListRsp{
				StatusCode: -1,
				StatusMsg: "获取视频失败",
				VideoList: nil,
			})
		}
		videoList = append(videoList, video)
	}

	// 获取视频详细信息
	videoInfos := make([]serializer.Video, 0)
	for _, video := range videoList {

		var author serializer.User
		// 通过AuthorId查找视频作者相关信息
		model.DB.Model(&serializer.User{}).Where("id = ?", video.AuthorId).First(&author)

		// 作者相关信息

		var isFollowCount int64
		if user_id > 0{
			// 查询是否已follow
			model.DB.Model(&model.Follow{}).Where("to_user_id = ? and user_id = ?", author.Id, userId).Count(&isFollowCount)
			author.IsFollow = isFollowCount > 0
		}

		videoInfos = append(videoInfos, serializer.Video{
			Id:            int64(video.ID),
			Author:        author,
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    true,
			Title:         video.Title,
		})
	}

	if
	StatusCode == 0{
		c.JSON(200, FavoriteListRsp{
			StatusCode: 0,
			StatusMsg: "操作成功",
			VideoList: videoInfos,
		})
	}


}
