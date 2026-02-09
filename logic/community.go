package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"

	"github.com/gin-gonic/gin"
)

func GetCommunityList(c *gin.Context) ([]*models.Community, error) {
	return mysql.GetCommunityList()
}

func GetCommunityDetail(c *gin.Context, id int64) (*models.CommunityDetail, error) {
	return mysql.GetCommunityDetail(id)
}
