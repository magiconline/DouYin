package cache

//import (
//	_"github.com/go-sql-driver/mysql"
//	"github.com/jmoiron/sqlx"
//	"DouYin/repository"
//	"github.com/go-redis/redis/v9"
//)
//type User struct {
//	UserId        int64  `gorm:"column:user_id; primary_key; AUTO_INCREMENT"`
//	UserName      string `gorm:"column:user_name"`
//	FollowCount   int    `gorm:"column:follow_count"`
//	FollowerCount int    `gorm:"column:follower_count"`
//}
//
//func CacheUser2Redis(u []User) error{
//	//1.连接mysql数据库
//	db, _ := sqlx.Connect("mysql","root:123456@tcp(127.0.0.1:3306)/douyin")
//	defer db.Close()
//	//2.选择要缓存的字段
//	var userinfo []User
//	err := db.Select(&u, "select user_id, user_name, follow_count, fllower_count from user")
//
//	//3.错误处理
//	if err != nil {
//		return err
//	}
//	//4.将数据写入redis
//}
