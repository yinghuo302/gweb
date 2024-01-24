package router

type BasicRouter struct {
	handler_mp map[string]any
}

func (router *BasicRouter) AddRoute(method string, pattern string, handlers any) {
	router.handler_mp[method+"-"+pattern] = handlers
}

func (router *BasicRouter) GetHandlers(method string, pattern string) (map[string]string, any) {
	handlers, ok := router.handler_mp[method+"-"+pattern]
	if !ok {
		handlers = nil
	}
	return nil, handlers
}
