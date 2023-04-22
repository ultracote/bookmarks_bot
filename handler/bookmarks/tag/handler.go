package tag

import (
	"bookmarks/db/entity"
	"bookmarks/routes"
	"bookmarks/service/router"
	"errors"
	"fmt"
	"github.com/mymmrac/telego"
	"gorm.io/gorm"
	"strconv"
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
	if len(request.GetArgs()) < 1 {
		return nil, errors.New("invalid arguments count")
	}

	// get bookmark id from args
	bookmarkID := request.GetArg(0)
	// get tag id from args
	tagID := request.GetArg(1)

	// convert bookmark id to int64
	bookmarkIDInt64, err := strconv.ParseInt(bookmarkID, 10, 64)
	if err != nil {
		return nil, errors.New("invalid bookmark id")
	}

	if len(tagID) != 0 {
		// convert tag id to int64
		tagIDInt64, err := strconv.ParseInt(tagID, 10, 64)
		if err != nil {
			return nil, errors.New("invalid tag id")
		}

		// save bookmark to tag
		err = h.addTagToBookmark(bookmarkIDInt64, tagIDInt64)
		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, fmt.Errorf("error updating user state: %w", err)
	}

	// retrieve all tags sort by linked to bookmark
	var tags []struct {
		Title      string
		TagID      int64
		RelationID *int64
	}
	err = h.db.
		Table("tags").
		Select("tags.title, tags.id as tag_id, bookmark_tags.tag_id as relation_id").
		Joins("LEFT JOIN bookmark_tags ON bookmark_tags.tag_id = tags.id AND bookmark_tags.bookmark_id = ?", bookmarkIDInt64).
		//Where("", ).
		Where("tags.user_id = ?", request.GetChatID()).
		Order("bookmark_tags.bookmark_id DESC").
		Find(&tags).Error

	if err != nil {
		return nil, fmt.Errorf("error retrieving tags: %w", err)
	}

	// send tags to user as keyboard
	var keyboard [][]telego.InlineKeyboardButton
	for _, tag := range tags {
		isLinked := tag.RelationID != nil
		text := tag.Title
		route := routes.BookmarkTag
		if isLinked {
			text = "✅" + text
			route = routes.BookmarkUntag
		}

		keyboard = append(keyboard, []telego.InlineKeyboardButton{
			{
				Text:         text,
				CallbackData: fmt.Sprintf(route+"?%s.%d", bookmarkID, tag.TagID),
			},
		})
	}

	// add button to create new tag
	keyboard = append(keyboard, []telego.InlineKeyboardButton{
		{
			Text:         "➕Новый тег",
			CallbackData: routes.TagCreateWait,
		},
	})

	// add button to create new tag
	keyboard = append(keyboard, []telego.InlineKeyboardButton{
		{
			Text:         "🆗Готово",
			CallbackData: routes.TagCreateDone,
		},
	})

	// send message to user
	return router.NewResponse().AddMessage(&telego.SendMessageParams{
		Text: "Select tag",
		ReplyMarkup: &telego.InlineKeyboardMarkup{
			InlineKeyboard: keyboard,
		},
	}), nil
}

func (h *Handler) addTagToBookmark(bookmarkID, tagID int64) error {
	bookmarkTag := entity.BookmarkTag{
		BookmarkID: bookmarkID,
		TagID:      tagID,
	}
	err := h.db.Create(&bookmarkTag).Error
	if err != nil {
		return fmt.Errorf("error saving bookmark to tag: %w", err)
	}

	return nil
}
