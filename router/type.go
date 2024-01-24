package router

type IRouter interface {
	AddRoute(method string, pattern string, handlers any)
	GetHandlers(method string, pattern string) (map[string]string, any)
}

func NewBasicRouter() IRouter {
	return &BasicRouter{handler_mp: make(map[string]any)}
}

func NewTrieRouter() IRouter {
	return &TrieRouter{childs: map[string]*TrieRouter{}}
}
