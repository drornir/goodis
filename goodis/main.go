package goodis

import "github.com/drornir/goodis/goodis/request"

type App interface {
	Handle(string) (string, error)
}

type app struct {
}

func New() App {
	return new(app)
}

func (a *app) Handle(req string) (string, error) {
	handler := new(request.Handler)
	return handler.Handle(req)
}
