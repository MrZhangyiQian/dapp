package main

import (
	"bolg/database"
	"bolg/handlers"
	"bolg/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	database.InitDB()
	defer database.DB.Close()

	r := gin.Default()

	// 用户相关路由
	r.POST("/register", handlers.RegisterUser(database.DB))
	r.POST("/login", handlers.LoginUser(database.DB))

	// 文章相关路由
	postGroup := r.Group("/posts")
	{
		postGroup.POST("", middleware.AuthMiddleware(), handlers.CreatePost)
		postGroup.GET("", handlers.GetPosts)
		postGroup.GET("/:id", handlers.GetPost)
		postGroup.PUT("/:id", middleware.AuthMiddleware(), handlers.UpdatePost)
		postGroup.DELETE("/:id", middleware.AuthMiddleware(), handlers.DeletePost)
	}

	// 评论相关路由
	commentGroup := r.Group("/comments")
	{
		commentGroup.POST("/:post_id", middleware.AuthMiddleware(), handlers.CreateComment)
		commentGroup.GET("/:post_id", handlers.GetCommentsByPost)
	}

	r.Run(":8080")
}
