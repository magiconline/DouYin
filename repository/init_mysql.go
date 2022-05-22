package repository

import (
	"fmt"
	mysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	mysqlUrl = "root:wang186017@tcp(127.0.0.1:3306)/douyin"
	DB       *gorm.DB
)

func Init() error {
	var err error
	DB, err = gorm.Open(mysql.Open(mysqlUrl), &gorm.Config{})
	DB.AutoMigrate(&Star{})
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
