package list

import (
	"bookmarks/db/entity"
	"bookmarks/service/router"
	"fmt"
	"github.com/mymmrac/telego"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{
		db: db,
	}
}

func (h *Handler) Handle(request *router.Request) (*router.Response, error) {
	// get bookmarks from db
	var bookmarks []entity.Bookmark
	err := h.db.Where("user_id = ?", request.GetChatID()).Find(&bookmarks).Error
	if err != nil {
		return nil, fmt.Errorf("error while getting bookmarks: %w", err)
	}

	// counter
	var i int

	var text string
	for _, bookmark := range bookmarks {
		i++
		text += fmt.Sprintf("%d. %s\n\n", i, bookmark.Text)
	}

	// if no bookmarks
	if i == 0 {
		text = "Нет закладок"
	}

	return router.NewResponse().AddMessage(&telego.SendMessageParams{
		Text: text,
	}), nil
}
