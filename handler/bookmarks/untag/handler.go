package untag

import (
	"bookmarks/db/entity"
	"bookmarks/routes"
	"bookmarks/service/router"
	"errors"
	"fmt"
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
	if len(request.GetArgs()) < 2 {
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

	// convert tag id to int64
	tagIDInt64, err := strconv.ParseInt(tagID, 10, 64)
	if err != nil {
		return nil, errors.New("invalid tag id")
	}

	// save bookmark to tag
	err = h.removeTagFromBookmark(bookmarkIDInt64, tagIDInt64)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, fmt.Errorf("error updating user state: %w", err)
	}

	// send message to user
	return router.NewResponse().SetRedirect(fmt.Sprintf(routes.BookmarkTag+"?%s", bookmarkID)), nil
}

func (h *Handler) removeTagFromBookmark(bookmarkID, tagID int64) error {
	// remove tag from bookmark
	err := h.db.
		Where("bookmark_id = ? AND tag_id = ?", bookmarkID, tagID).
		Delete(&entity.BookmarkTag{}).Error
	if err != nil {
		return fmt.Errorf("error removing bookmark to tag: %w", err)
	}

	return nil
}
