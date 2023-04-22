package tag

import (
	"bookmarks/db/entity"
	"bookmarks/service/router"
	"errors"
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
	// check if there are 1 argument
	if len(request.GetArgs()) < 1 {
		return nil, errors.New("invalid arguments count")
	}

	// get tag id from args
	tagID := request.GetArg(0)

	// convert tag id to int64
	tagIDInt64, err := strconv.ParseInt(tagID, 10, 64)
	if err != nil {
		return nil, errors.New("invalid tag id")
	}

	// retrieve bookmarks by tag id
	var bookmarks []entity.Bookmark
	tx := h.db.Where("id IN (SELECT bookmark_id FROM bookmark_tags WHERE tag_id = ?)", tagIDInt64).Find(&bookmarks)
	if tx.Error != nil {
		return nil, errors.New("error retrieving bookmarks")
	}

	// send message if there are no bookmarks in tag
	if len(bookmarks) == 0 {
		return router.NewResponse().AddMessage(&telego.SendMessageParams{
			Text: "There are no bookmarks in tag", // todo
		}), nil
	}

	// send bookmarks to user
	var text string
	for index, bookmark := range bookmarks {
		text += strconv.Itoa(index+1) + ". " + bookmark.Text + "\n\n"
	}

	return router.NewResponse().AddMessage(&telego.SendMessageParams{
		Text: text,
	}), nil
}
