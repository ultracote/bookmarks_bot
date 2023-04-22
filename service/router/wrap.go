package router

func Wrap(handleFunc HandleFunc) Handler {
	return &wrapHandler{
		handleFunc: handleFunc,
	}
}

type wrapHandler struct {
	handleFunc HandleFunc
}

func (w *wrapHandler) Handle(request *Request) (*Response, error) {
	return w.handleFunc(request)
}
