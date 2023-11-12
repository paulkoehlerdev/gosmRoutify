package router

type Application interface {
}

type impl struct {
}

func New() Application {
	return &impl{}
}
