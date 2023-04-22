package redirect

import (
	"bookmarks/service/router"
	"errors"
)

type Handler struct {
	handler router.Handler
}

func NewHandler(handler router.Handler) *Handler {
	return &Handler{handler: handler}
}

func (h *Handler) Handle(request *router.Request) (*router.Response, error) {
	if request.RedirectsCount > 3 {
		return nil, errors.New("too many redirects")
	}

	response, err := h.handler.Handle(request)
	if err != nil {
		return nil, err
	}

	if response.HasRedirect() {
		request.RedirectsCount++
		redirectDestination := response.GetRedirect().Route
		request.SetRoute(&redirectDestination)

		return h.Handle(request)
	} else {
		// clear route after redirect
		request.SetRoute(nil)
	}

	return response, nil
}
