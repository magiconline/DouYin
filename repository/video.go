package repository

type VideoTable struct {
	VideoId       int64
	UserId        int64
	PlayUrl       string
	CoverUrl      string
	FavoriteCount int32
	CommentCount  int32
	UploadTime    int32
}
