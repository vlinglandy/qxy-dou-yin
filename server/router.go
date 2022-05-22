package server

import (
	"os"
	"qxy-dy/api"
	"qxy-dy/middleware"

	"github.com/gin-gonic/gin"
)

// NewRouter 路由配置
func NewRouter() *gin.Engine {
	r := gin.Default()

	// 中间件, 顺序不能改
	r.Use(middleware.Session(os.Getenv("SESSION_SECRET")))
	r.Use(middleware.Cors())
	r.Use(middleware.CurrentUser())

	// 路由
	douyin := r.Group("/douyin")
	{
		// 一个简单的示例，可供大家学习
		douyin.POST("ping", api.Ping)

		// 用户注册
		douyin.POST("user/register", api.UserRegister)

		// 用户登录
		douyin.POST("user/login", api.UserLogin)

		// 视频流接口
		douyin.GET("feed", api.GetVideo)

		// 需要登录保护的，说白了就是不能用不能用的接口，视频流不登陆也能用
		auth := douyin.Group("")
		auth.Use(middleware.AuthRequired())
		{
			// auth.GET("user",api)

			// 视频投稿
			auth.POST("publish/action", api.PublishAction)

			// 发布视频列表
			auth.GET("publish/list", api.PubishList)

			// 点赞操作
			auth.POST("favorite/action", api.FavoriteAction)

			// 点赞列表
			auth.GET("favorite/list", api.FavoriteList)

			// 评论操作
			auth.POST("comment/action", api.ComentAction)

			// 视频评论列表
			auth.GET("comment/list", api.CommentList)

			// 关系操作
			auth.POST("relation/action", api.RelationAction)

			// 用户关注列表
			auth.GET("relation/follow/list", api.UserFollowList)

			// 用户粉丝列表
			auth.GET("relation/follower/list", api.UserLoverList)
		}
	}
	return r
}
