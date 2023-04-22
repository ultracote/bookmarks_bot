package wait

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
	// set user state to tag create
	state := entity.UserState{
		UserID:           request.GetChatID(),
		NextMessageRoute: routes.TagCreate,
	}

	err := h.db.Save(&state).Error
	if err != nil {
		return nil, fmt.Errorf("error while saving user state: %w", err)
	}

	return router.NewResponse().AddMessage(&telego.SendMessageParams{
		Text: "Введите название тега ответом на это сообщение",
	}), nil
}
