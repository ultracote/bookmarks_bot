package entity

type BookmarkTag struct {
	BookmarkID int64 `gorm:"primarykey"`
	TagID      int64 `gorm:"primarykey"`
}
