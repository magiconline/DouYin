package service

import "DouYin/repository"

func FavoriteOperation(star *repository.Star) {
	repository.NewStarDaoInstance().FavorableOperation(star)
}

func QueryByUserIdAndVideoId(userId, videoId uint64) *repository.Star {
	return repository.NewStarDaoInstance().QueryByUserIdAndVideoId(userId, videoId)
}
