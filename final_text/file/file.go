// Package file 提供文件管理相关的功能，包括上传、删除、恢复、收藏和取消收藏文件等。
package file

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type File struct {
	ID        uint64    `gorm:"primarykey"`
	UserID    uint64    `gorm:"index"`
	Filename  string    `gorm:"size:255"`
	Filepath  string    `gorm:"size:512"`
	Filesize  int64
	SharedKey string
	DeletedAt *time.Time
	CreatedAt time.Time
}

type favoriteFile struct {
	ID     uint64 `gorm:"primarykey"`
	UserID uint64 `gorm:"index"`
	FileID uint64 `gorm:"index"`
}

var db *gorm.DB

func InitDB(database *gorm.DB) {
	db = database
	db.AutoMigrate(&File{}, &favoriteFile{})
}

// UploadFile 处理文件上传请求
func UploadFile(c *gin.Context) {
	userID := c.MustGet("user_id").(uint64)
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "err.Error()",
		})
		return
	}
	uploadPath := fmt.Sprintf("./uploads/%d/", userID)
	os.MkdirAll(uploadPath, os.ModePerm)
	filePath := filepath.Join(uploadPath, file.Filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "err.Error()",
		})
		return
	}
	newFile := File{
		UserID:   userID,
		Filename: file.Filename,
		Filepath: filePath,
		Filesize: file.Size,
	}
	db.Create(&newFile)
	c.JSON(http.StatusOK, gin.H{
		"message": "upload success","file_id": newFile.ID,
	})
}

//DeleteFile 删除用户的文件
func DeleteFile(c *gin.Context) {
	userID := c.MustGet("user_id").(uint64)
	fileID := c.Query("file_id")
	var file File
	if err := db.Where("id = ? AND user_id = ?", fileID, userID).First(&file).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	if err := os.Remove(file.Filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file from storage"})
		return
	}
	now:= time.Now()
	db.Model(&file).Update("deleted_at", &now)
	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}

//RestoreFile 恢复用户的文件
func RestoreFile(c *gin.Context) {
	userID := c.MustGet("user_id").(uint64)
	fileID := c.Param("file_id")
	var file File
	if err := db.Where("id = ? AND user_id = ?", fileID, userID).First(&file).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	db.Model(&file).Update("deleted_at", nil)
	c.JSON(http.StatusOK, gin.H{"message": "File restored successfully"})
}

//FavoriteFile 收藏用户的文件
func FavoriteFile(c *gin.Context) {
	userID := c.MustGet("user_id").(uint64)
	fileIDStr := c.Query("file_id")
	fileID, err := strconv.ParseUint(fileIDStr, 10, 64)
if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file_id"})
    return
}
	db.Create(&favoriteFile{
		UserID: userID,
		FileID: fileID,
	})
	c.JSON(http.StatusOK, gin.H{"message": "File favorited successfully"})
}

func stringToUint64(s string) uint64 {
	var result uint64
	fmt.Sscanf(s, "%d", &result)
	return result
}

//UnfavoriteFile 取消收藏用户的文件
func UnfavoriteFile(c *gin.Context) {
	userID := c.MustGet("user_id").(uint64)
	fileIDStr := c.Query("file_id")
	fileID, err := strconv.ParseUint(fileIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file_id"})
		return
	}
	db.Where("user_id = ? AND file_id = ?", userID, fileID).Delete(&favoriteFile{})
	c.JSON(http.StatusOK, gin.H{"message": "File unfavorited successfully"})
}

//ListFavoriteFiles 列出用户收藏的文件
func ListFavoriteFiles(c *gin.Context) {
	userID := c.MustGet("user_id").(uint64)
	var favFiles []favoriteFile
	db.Where("user_id = ?", userID).Find(&favFiles)
	var fileIDs []uint64
	for _, fav := range favFiles {
		fileIDs = append(fileIDs, fav.FileID)
	}
	var files []File
	db.Where("id IN ?", fileIDs).Find(&files)
	c.JSON(http.StatusOK, gin.H{"favorite_files": files})
}