package repository

type StarTable struct {
	ID         uint64
	UserId     uint64
	StarType   uint8
	StarTypeId uint64
	Status     uint8
}

func (v StarTable) TableName() string {
	return "star"
}
