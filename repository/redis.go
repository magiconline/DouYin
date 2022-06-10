package repository

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
)

type UserInfomation struct {
	UserId        int64  `gorm:"column:user_id; primary_key; AUTO_INCREMENT"`
	UserName      string `gorm:"column:user_name"`
	FollowCount   int    `gorm:"column:follow_count"`
	FollowerCount int    `gorm:"column:follower_count"`
}

// func GetValueByKey(key string) (*UserInfomation, error) {
// 	var user UserInfomation
// 	user := RDB.Get(key)
// 	return &user, err
// }
func FindUserFromRedis(c *gin.Context) *redis.StringCmd {
	userID := c.Query("user_id")

	return RDB.Get(c, userID)

}
