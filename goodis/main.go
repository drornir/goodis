package goodis

type App interface {
	Handler() Handler
}

type app struct {
	handler Handler
}

func New() App {
	return &app{
		handler: AppHandler(),
	}
}

func (a app) Handler() Handler {
	return a.handler
}
