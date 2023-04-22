package tags

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
	// reply with all user tags as keyboard
	var tags []entity.Tag
	tx := h.db.Where("user_id = ?", request.GetChatID()).Find(&tags)
	if tx.Error != nil {
		return nil, fmt.Errorf("error retrieving tags: %w", tx.Error)
	}

	// create keyboard
	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: make([][]telego.InlineKeyboardButton, 0),
	}

	// add button for each tag
	for _, tag := range tags {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []telego.InlineKeyboardButton{
			{
				Text:         tag.Title,
				CallbackData: fmt.Sprintf(routes.TagShow+"?%d", tag.ID),
			},
		})
	}

	// send message with keyboard if there are tags or message without keyboard if there are no tags
	if len(tags) > 0 {
		return router.NewResponse().AddMessage(&telego.SendMessageParams{
			Text:        "Select tag",
			ReplyMarkup: keyboard,
		}), nil
	} else {
		return router.NewResponse().AddMessage(&telego.SendMessageParams{
			Text: "You have no tags",
		}), nil
	}
}
