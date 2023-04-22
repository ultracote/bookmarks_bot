package entity

type UserState struct {
	UserID int64 `gorm:"primarykey;autoIncrement:false"`
	//LastMessageID    int64  `gorm:"index;not null"`
	NextMessageRoute string `gorm:"type:text"`
}
