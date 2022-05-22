package service

import "DouYin/repository"

func AddStar(userId, videoId uint64) {
	repository.NewStarDaoInstance().AddStar(userId, videoId)
}

func DeleteStar(userId, videoId uint64) {
	repository.NewStarDaoInstance().DeleteStar(userId, videoId)
}
