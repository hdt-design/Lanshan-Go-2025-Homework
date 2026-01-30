package main

import (
	"final_text/file"
	"final_text/user"

	"github.com/gin-gonic/gin"
)

func main() {
	user.InitDB()
	file.InitDB(user.DB()) // 把 user 包的 DB 传给 file 包

	r := gin.Default()
	// 文件模块
	// 用户模块OST("/register", user.Register)
	r.POST("/login", user.Login)

	auth := r.Group("/user")
	auth.Use(user.AuthMiddleware())
	{
		auth.GET("/info", user.GetUserInfo)
		auth.PUT("/update", user.UpdateUserInfo)

		// 文件操作
		auth.POST("/upload", file.UploadFile)
		auth.DELETE("/delete", file.DeleteFile)
		auth.PUT("/restore", file.RestoreFile)
		auth.POST("/favorite", file.FavoriteFile)
	}

	r.Run(":8080")
}
