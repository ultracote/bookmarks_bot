package start

import (
	"bookmarks/service/router"
	"github.com/mymmrac/telego"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(request *router.Request) (*router.Response, error) {
	return router.NewResponse().AddMessage(&telego.SendMessageParams{
		Text: "Hello, I'm bookmarks bot!",
	}), nil
}
