/*
 * @Author: zanilia
 * @Date: 2022-10-02 19:32:03
 * @LastEditTime: 2022-10-05 11:21:03
 * @Descripttion:
 */
package router

import (
	"strings"
)

type TrieRouter struct {
	childs   map[string]*TrieRouter
	Param    string
	handlers any
}

func (root *TrieRouter) AddRoute(method string, pattern string, handlers any) {
	parts := strings.Split(pattern, "/")
	node := root
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}
		param, childName := "", part
		if part[0] == ':' || part[0] == '*' {
			param, childName = part, ""
		}
		if _, ok := node.childs[childName]; !ok {
			node.childs[childName] = &TrieRouter{
				childs: make(map[string]*TrieRouter),
			}
		}
		node = node.childs[childName]
		node.Param = param
		if part[0] == '*' {
			break
		}
	}
	node.childs[method] = &TrieRouter{
		childs:   make(map[string]*TrieRouter),
		handlers: handlers,
	}
}

func (router *TrieRouter) GetHandlers(method string, pattern string) (map[string]string, any) {
	parts := strings.Split(pattern, "/")
	node := router
	params := make(map[string]string)
	var ok bool
	for idx, part := range parts {
		if len(part) == 0 {
			continue
		}
		node, ok = getNode(node.childs, part, "")
		if !ok {
			return nil, nil
		}
		if len(node.Param) != 0 && node.Param[0] == '*' {
			params[node.Param[1:]] = strings.Join(parts[idx:], "/")
			break
		} else if len(node.Param) != 0 && node.Param[0] == ':' {
			params[node.Param[1:]] = part
		}
	}
	node, ok = getNode(node.childs, method, "ALL")
	if !ok {
		return nil, nil
	} else {
		return params, node.handlers
	}
}

func getNode(mp map[string]*TrieRouter, first, second string) (*TrieRouter, bool) {
	node, ok := mp[first]
	if !ok {
		node, ok = mp[second]
	}
	return node, ok
}
