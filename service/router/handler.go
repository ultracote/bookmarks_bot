package router

import (
	"errors"
	"github.com/mymmrac/telego"
)

type Handler interface {
	Handle(request *Request) (*Response, error)
}

type HandleFunc func(request *Request) (*Response, error)

type Request struct {
	Update         *telego.Update
	Args           []string
	ChatID         int64
	Route          *string
	RedirectsCount int
}

func NewRequest(update *telego.Update, chatID int64, args ...string) *Request {
	return &Request{
		Update: update,
		ChatID: chatID,
		Args:   args,
	}
}

func (r *Request) GetArg(index int) string {
	if len(r.Args) > index {
		return r.Args[index]
	}
	return ""
}

//func (r *Request) AddArgs(args ...string) {
//	r.Args = append(r.Args, args...)
//}

func (r *Request) SetArgs(args ...string) {
	r.Args = args
}

func (r *Request) RequireArg(index int) (string, error) {
	if len(r.Args) > index {
		return r.Args[index], nil
	}
	return "", errors.New("argument not found")
}

func (r *Request) GetArgs() []string {
	return r.Args
}

func (r *Request) GetUpdate() *telego.Update {
	return r.Update
}

func (r *Request) GetChatID() int64 {
	return r.ChatID
}

func (r *Request) SetRoute(route *string) {
	r.Route = route
}

type Response struct {
	Messages []*telego.SendMessageParams
	Redirect *Redirect
}

func NewResponse() *Response {
	return &Response{}
}

func (r *Response) AddMessage(message ...*telego.SendMessageParams) *Response {
	r.Messages = append(r.Messages, message...)
	return r
}

func (r *Response) AddMessages(messages []*telego.SendMessageParams) *Response {
	r.Messages = append(r.Messages, messages...)
	return r
}

func (r *Response) SetMessages(messages []*telego.SendMessageParams) *Response {
	r.Messages = messages
	return r
}

func (r *Response) SetRedirect(route string) *Response {
	r.Redirect = &Redirect{
		Route: route,
	}
	return r
}

func (r *Response) GetMessages() []*telego.SendMessageParams {
	return r.Messages
}

func (r *Response) GetRedirect() *Redirect {
	return r.Redirect
}

func (r *Response) HasRedirect() bool {
	return r.Redirect != nil
}

type Redirect struct {
	Route string
}
