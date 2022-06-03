package repository

type Relation struct {
	UserID   uint64 `gorm:"primaryKey;notNULL"` // 用户
	ToUserID uint64 `gorm:"notNULL;index"`      // 被关注的用户
	Relation bool   `gorm:"notNULL"`
}

// 关注操作
// action=true表示关注
// action=false表示取消关注
func Action(userID uint64, toUserID uint64, action bool) error {
	if action {
		err := DB.FirstOrCreate(&Relation{UserID: userID, ToUserID: toUserID}, Relation{UserID: userID, ToUserID: toUserID, Relation: true}).Error
		return err
	} else {
		err := DB.Model(&Relation{UserID: userID, ToUserID: toUserID}).Update("relation", false).Error
		return err
	}
}

// 关注列表
func FollowList(userID uint64) (*[]Relation, error) {
	var results []Relation
	err := DB.Where(&Relation{UserID: userID, Relation: true}).Select("to_user_id").Find(&results).Error

	return &results, err
}

// 粉丝列表
func FollowerList(userID uint64) (*[]Relation, error) {
	var results []Relation
	err := DB.Where(&Relation{ToUserID: userID, Relation: true}).Select("user_id").Find(&results).Error

	return &results, err
}

// 判断userID是否关注了toUserID
func IsFollower(userID uint64, toUserID uint64) (bool, error) {
	relation := &Relation{}
	err := DB.Where(&Relation{UserID: userID, ToUserID: toUserID, Relation: true}).Limit(1).Find(&relation).Error
	return relation.Relation, err
}

//
