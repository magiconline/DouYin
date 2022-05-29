package repository

// "DouYin/json"

type User struct {
	// gorm.Model
	UserId        int64  `gorm:"column:user_id; primary_key; AUTO_INCREMENT"`
	Email         string `gorm:"column:email"`
	UserName      string `gorm:"column:user_name"`
	Password      string `gorm:"column:password"`
	FollowCount   int    `gorm:"column:follow_count"`
	FollowerCount int    `gorm:"column:follower_count"`
}

//根据email pwd查找用户
func FindUserbyEmailandPwd(email, pwd string) (*User, error) {
	var user User
	err := DB.Table("user").Where("email = ? and password = ?", email, pwd).First(&user).Error
	return &user, err
}

//根据email查找用户
func FindUserbyEmail(email string) (*User, error) {
	var user User
	err := DB.Table("user").Where("email = ? ", email).First(&user).Error
	return &user, err
}

//根据userid查找用户
func FindUserbyID(userid string) (*User, error) {
	var user User
	err := DB.Table("user").Where("user_id = ?", userid).First(&user).Error
	return &user, err
}

// 创建用户
func CreateUser(email, pwd string) (*User, error) {
	// _, err1 := FindUserbyEmail(email)
	// if err1 == nil {
	user := User{
		Email:    email,
		UserName: "user_" + email,
		Password: pwd,
	}
	err := DB.Table("user").Create(&user).Error
	return &user, err
}
