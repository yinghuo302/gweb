package router

import (
	"testing"
)

type RouteInfo struct {
	Pattern  string
	Method   string
	Params   map[string]string
	Handlers int
}

func TestTrieRoute(t *testing.T) {
	trie := NewTrieRouter()
	routes := []RouteInfo{{
		Pattern:  "/api/test1/1",
		Method:   "POST",
		Handlers: 1,
	}, {
		Pattern:  "/api/test2/*filepath",
		Method:   "GET",
		Handlers: 2,
	}, {
		Pattern:  "/api/test3/:id/1/:arg",
		Method:   "GET",
		Handlers: 3,
	}, {
		Pattern:  "/api/test4/:id",
		Method:   "ALL",
		Handlers: 4,
	}}

	testRoute := []RouteInfo{{
		Pattern:  "/api/test1/1",
		Method:   "POST",
		Handlers: 1,
	}, {
		Pattern:  "/api/test1/1",
		Method:   "GET",
		Handlers: -1,
	}, {
		Pattern:  "/api/test1/2",
		Method:   "POST",
		Handlers: -1,
	}, {
		Pattern:  "/api/test2/img/1.jpg",
		Method:   "GET",
		Handlers: 2,
		Params:   map[string]string{"filepath": "img/1.jpg"},
	}, {
		Pattern:  "/api/test2/var/1.log",
		Method:   "GET",
		Handlers: 2,
		Params:   map[string]string{"filepath": "var/1.log"},
	}, {
		Pattern:  "/api/test2/var/1.log",
		Method:   "DELETE",
		Handlers: -1,
		Params:   nil,
	}, {
		Pattern:  "/api/test3/1/1/3",
		Method:   "GET",
		Handlers: 3,
		Params:   map[string]string{"id": "1", "arg": "3"},
	}, {
		Pattern:  "/api/test3/2/test/1/3",
		Method:   "GET",
		Handlers: -1,
		Params:   nil,
	}, {
		Pattern:  "/api/test3/1/2/3",
		Method:   "GET",
		Handlers: -1,
		Params:   nil,
	}, {
		Pattern:  "/api/test4/1",
		Method:   "POST",
		Handlers: 4,
		Params:   map[string]string{"id": "1"},
	}, {
		Pattern:  "/api/test4/2",
		Method:   "GET",
		Handlers: 4,
		Params:   map[string]string{"id": "2"},
	}, {
		Pattern:  "/api/test4",
		Method:   "GET",
		Handlers: -1,
		Params:   nil,
	}}

	for _, route := range routes {
		trie.AddRoute(route.Method, route.Pattern, route.Handlers)
	}

	for idx, testcase := range testRoute {
		params, handlers := trie.GetHandlers(testcase.Method, testcase.Pattern)
		if testcase.Handlers != handlersToInt(handlers) || !mapEqual(params, testcase.Params) {
			t.Fatalf("testcase %d: method: \"%s\" path: \"%s\" should get params: %v handlers: %v but get params: %v, handlers: %v", idx, testcase.Method, testcase.Pattern, testcase.Params, testcase.Handlers, params, handlers)
		}
	}

}

func handlersToInt(handlers any) int {
	if handlers == nil {
		return -1
	}
	return handlers.(int)
}

func mapEqual(a, b map[string]string) bool {

	if len(a) != len(b) {
		return false
	}
	for key, val := range a {
		if b_val, ok := b[key]; !ok || b_val != val {
			return false
		}
	}
	return true
}
