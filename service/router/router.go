package router

import (
	"bookmarks/db/entity"
	"errors"
	"fmt"
	"github.com/mymmrac/telego"
	"gorm.io/gorm"
	"log"
	"strings"
)

type Router struct {
	db              *gorm.DB
	fallbackHandler Handler
	errorHandler    Handler

	callbackHandlers map[string]Handler
	commandHandlers  map[string]Handler
}

func NewRouter(db *gorm.DB, fallbackHandler Handler, errorHandler Handler) *Router {
	return &Router{
		db:              db,
		fallbackHandler: fallbackHandler,
		errorHandler:    errorHandler,

		callbackHandlers: make(map[string]Handler),
		commandHandlers:  make(map[string]Handler),
	}
}

func (r *Router) Handle(request *Request) (*Response, error) {
	var handler Handler

	// todo странная ситуация с аргументами

	update := request.GetUpdate()
	var callbackRoute string

	if request.Route != nil {
		callbackRoute = *request.Route
	} else if update.CallbackQuery != nil {
		callbackRoute = update.CallbackQuery.Data
	} else {
		// try to find callback handler if user has next message route in state
		// retrieve user state
		chatID := r.getChatID(update)

		var userState entity.UserState
		err := r.db.Where("user_id = ?", chatID).FirstOrInit(&userState).Error // todo стремная конструкция
		if err != nil {
			return nil, fmt.Errorf("error while getting user state: %w", err)
		}

		// if user has next message route
		if userState.NextMessageRoute != "" {
			// todo clean next message route here
			callbackRoute = userState.NextMessageRoute

			userState.NextMessageRoute = ""
			err = r.db.Save(&userState).Error
			if err != nil {
				return nil, fmt.Errorf("error while saving user state in router: %w", err)
			}
		}
	}

	if callbackRoute != "" {
		// split data to command and arguments by ?
		parts := strings.SplitN(callbackRoute, "?", 2)
		route := parts[0]

		if callbackHandler, ok := r.callbackHandlers[route]; ok {
			handler = callbackHandler
			var argsString string
			if len(parts) > 1 {
				argsString = parts[1]
			}

			// split args to list of arguments by .
			request.SetArgs(strings.Split(argsString, ".")...)
		}
	}

	// todo check if its command
	if handler == nil && update.Message != nil {
		text := update.Message.Text
		// split text to command and value by space
		parts := strings.SplitN(text, " ", 2)
		cmd := parts[0]

		if commandHandler, ok := r.commandHandlers[cmd]; ok {
			handler = commandHandler

			var value string
			if len(parts) > 1 {
				value = parts[1]
			}

			request.SetArgs(value)
		}
	}

	if handler == nil {
		handler = r.fallbackHandler
	}

	if handler == nil {
		return nil, errors.New("no handler found")
	}

	response, err := handler.Handle(request)
	if err != nil {
		return r.HandleError(request, err)
	}

	return response, nil
}

func (r *Router) AddCommandHandler(command string, handler Handler) {
	r.commandHandlers[command] = handler
}

func (r *Router) AddCallbackHandler(callback string, handler Handler) {
	r.callbackHandlers[callback] = handler
}

func (r *Router) SetErrorHandler(handler Handler) {
	r.errorHandler = handler
}

func (r *Router) HandleError(request *Request, err error) (*Response, error) {
	// log error
	log.Printf("Error caught: %s\n", err)
	// todo improve error handling

	if r.errorHandler == nil {
		return nil, fmt.Errorf("uncaught error: %w", err)
	}

	return r.errorHandler.Handle(request)
}

func (r *Router) getChatID(update *telego.Update) int64 {
	if update.Message != nil {
		return update.Message.Chat.ID
	}

	if update.CallbackQuery != nil {
		return update.CallbackQuery.Message.Chat.ID
	}

	return 0
}
