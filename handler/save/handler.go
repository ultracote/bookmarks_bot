package save

import (
	"bookmarks/db/entity"
	"bookmarks/routes"
	"bookmarks/service/router"
	"fmt"
	"github.com/mymmrac/telego"
	"gorm.io/gorm"
	"strings"
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
	// value is first argument if it exists or return error
	var value string
	if len(request.GetArgs()) > 0 {
		value = request.GetArg(0)
	} else {
		return nil, fmt.Errorf("no value provided")
	}

	value = strings.TrimSpace(value)

	if len(value) == 0 {
		return router.NewResponse().AddMessage(&telego.SendMessageParams{
			Text: "Закладка не может быть пустой",
		}), nil
	}

	// save bookmark to db
	bookmark := entity.Bookmark{
		UserID: request.GetChatID(),
		Text:   value,
	}
	tx := h.db.Create(&bookmark)
	if tx.Error != nil {
		return nil, fmt.Errorf("error while saving bookmark: %w", tx.Error)
	}

	// get id of inserted bookmark
	bookmarkID := bookmark.ID

	// retrieve user tags
	var tags []entity.Tag
	err := h.db.Where("user_id = ?", request.GetChatID()).Find(&tags).Error
	if err != nil {
		return nil, fmt.Errorf("error while getting tags: %w", err)
	}

	// create buttons for each tag
	var buttons []telego.InlineKeyboardButton
	for _, tag := range tags {
		buttons = append(buttons, telego.InlineKeyboardButton{
			Text:         tag.Title,
			CallbackData: fmt.Sprintf(routes.BookmarkTag+"?%d.%d", bookmarkID, tag.ID),
		})
	}

	// add button to create new tag
	buttons = append(buttons, telego.InlineKeyboardButton{
		Text:         "Create new tag",
		CallbackData: routes.TagCreateWait,
	})

	return router.NewResponse().AddMessage(&telego.SendMessageParams{
		Text: "Закладка сохранена. Выбери теги или создай новый",
		//ReplyMarkup: &telego.InlineKeyboardMarkup{
		//	InlineKeyboard: [][]telego.InlineKeyboardButton{buttons},
		//},
	}).SetRedirect(fmt.Sprintf(routes.BookmarkTag+"?%d", bookmarkID)), nil
}
