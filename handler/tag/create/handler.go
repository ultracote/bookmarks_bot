package tag

import (
	"bookmarks/db/entity"
	"bookmarks/routes"
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
	// if message empty return error
	if request.GetUpdate().Message == nil || request.GetUpdate().Message.Text == "" {
		// send message that tag can't be empty
		return router.NewResponse().AddMessage(&telego.SendMessageParams{
			Text: "Тег не может быть пустым",
		}), nil
	}

	// create tag
	tag := entity.Tag{
		UserID: request.GetChatID(),
		Title:  request.GetUpdate().Message.Text,
	}
	tx := h.db.Create(&tag)
	if tx.Error != nil {
		return nil, fmt.Errorf("error while creating tag: %w", tx.Error)
	}

	// todo receive tag id as argument
	// retrieve latest user bookmark
	var bookmark entity.Bookmark
	err := h.db.Where("user_id = ?", request.GetChatID()).Last(&bookmark).Error
	if err != nil {
		return nil, fmt.Errorf("error while getting latest bookmark: %w", err)
	}

	// create bookmark-tag relation
	bookmarkTag := entity.BookmarkTag{
		BookmarkID: bookmark.ID,
		TagID:      tag.ID,
	}

	tx = h.db.Create(&bookmarkTag)
	if tx.Error != nil {
		return nil, fmt.Errorf("error while creating bookmark-tag relation: %w", tx.Error)
	}

	return router.NewResponse().AddMessage(&telego.SendMessageParams{
		Text: "Тег создан: " + tag.Title + " и добавлен к закладке",
	}).SetRedirect(fmt.Sprintf(routes.BookmarkTag+"?%d", bookmark.ID)), nil
}
