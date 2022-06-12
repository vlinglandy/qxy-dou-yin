package api

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"qxy-dy/middleware"
	"qxy-dy/model"
	"qxy-dy/serializer"
	"qxy-dy/util"
)

type PublishService struct{}

func (PublishService) SaveVideo(authorId uint64, playUrl, coverUrl, title string) int64 {
	// 获取服务器的对外ip
	ip := util.GetIp()
	return model.DB.Create(&model.Video{
		AuthorId: uint(authorId),
		PlayUrl:  fmt.Sprintf("%s://%s:%s/%s", "http", ip, "8080", playUrl),
		CoverUrl: fmt.Sprintf("%s://%s:%s/%s", "http", ip, "8080", coverUrl),
		Title:    title,
	}).RowsAffected
}

func (PublishService) List(id string) []serializer.Video {
	var videos []model.Video
	model.DB.Model(&model.Video{}).Where("author_id = ?", id).Find(&videos)
	videoInfos := make([]serializer.Video, 0)
	for _, video := range videos {
		var author serializer.User
		// 通过AuthorId查找视频作者相关信息
		model.DB.Model(&serializer.User{}).Where("id = ?", video.AuthorId).First(&author)

		var isFollowCount, isFavoriteCount int64
		// 查询是否已follow
		model.DB.Model(&model.Follow{}).Where("to_user_id = ? and user_id = ?", author.Id, id).Count(&isFollowCount)
		author.IsFollow = isFollowCount > 0

		// 查询是否已喜欢
		model.DB.Model(&model.Favorite{}).Where("user_id = ? and video_id = ?", id, video.ID).Count(&isFavoriteCount)

		videoInfos = append(videoInfos, serializer.Video{
			Id:            int64(video.ID),
			Author:        author,
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    isFavoriteCount > 0,
			Title:         video.Title,
		})
	}
	return videoInfos
}

var publishService PublishService

// Publish 登录用户选择视频上传
func PublishAction(c *gin.Context) {
	// 获取参数token, 注意：这个token是存在text中的

	token, ok := c.GetPostForm("token")

	// 未获取到token
	if !ok {
		c.JSON(http.StatusOK, serializer.Response{
			StatusCode: 1,
			StatusMsg:  "未获取到token",
		})
		return
	}

	// 通过token获取id
	id, err := middleware.GetIdByToken(token)
	if err != nil {
		c.JSON(http.StatusOK, serializer.Response{
			StatusCode: 1,
			StatusMsg:  "由token获取id失败",
		})
		return
	}

	// 读取data
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, serializer.Response{
			StatusCode: 1,
			StatusMsg:  "加载data失败",
		})
		return
	}

	filename := filepath.Base(data.Filename)
	saveFile := fmt.Sprintf("public/%d_%s", id, filename)
	sf := []byte(saveFile)
	saveFileJepg := string(sf[:(len(sf)-4)]) + ".jpeg"

	title, _ := c.GetPostForm("title")

	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, serializer.Response{
			StatusCode: 1,
			StatusMsg:  "保存失败",
		})
		return
	} else {
		//截取某帧
		reader := ExampleReadFrameAsJpeg(saveFile, 1)
		img, err := imaging.Decode(reader)
		err = imaging.Save(img, saveFileJepg)
		if err != nil {
			c.JSON(http.StatusOK, serializer.Response{
				StatusCode: 1,
				StatusMsg:  "截图失败",
			})
			return
		}
		if row := publishService.SaveVideo(id, saveFile, saveFileJepg, title); row == 0 {
			c.JSON(http.StatusOK, serializer.Response{
				StatusCode: 1,
				StatusMsg:  "保存失败",
			})
			return
		}
	}

	// 上传成功
	c.JSON(http.StatusOK, serializer.Response{
		StatusCode: 0,
		StatusMsg:  saveFile + " uploaded successfully",
	})
}

// PublishList 登录用户所发表的所有视频
func PublishList(c *gin.Context) {
	// 获取参数user_id
	userId := c.Query("user_id")

	// 未获取到user_id
	if userId == "" {
		c.JSON(http.StatusOK, serializer.Response{
			StatusCode: 1,
			StatusMsg:  "未获取到user_id",
		})
		return
	}

	// 返回发表的视频
	c.JSON(http.StatusOK, serializer.PublishResponse{
		Response: serializer.Response{
			StatusCode: 0,
			StatusMsg:  "成功",
		},
		VideoInfoList: publishService.List(userId),
	})
}
func ExampleReadFrameAsJpeg(inFileName string, frameNum int) io.Reader {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(inFileName).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		panic(err)
	}
	return buf
}
