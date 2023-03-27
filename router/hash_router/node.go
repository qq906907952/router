package hash_router

import (
	"container/list"
	"github.com/qq906907952/router/router"
)

type node interface {
	route(routeKeys []string, idx int, routePath []string) (node, error)
	routeKey() string
	// the info must assure len at lease 3
	addRoute(id uint, info []string, idx int) error
	// return id
	deleteRoute(info []string, idx int) (uint, error)
}

type leafNode struct {
	id       uint
	endpoint string
}

func (l *leafNode) deleteRoute(info []string, idx int) (uint, error) {
	return l.id, nil
}

func (l *leafNode) addRoute(id uint, info []string, idx int) error {
	panic("unreachable")
}

func (l *leafNode) routeKey() string {
	return l.endpoint
}

func (l *leafNode) route(_ []string, _ int, _ []string) (node, error) {
	return l, nil
}

type routeNode struct {
	key string
	// the link list ele is node
	childMap map[string]*list.Element
	child    *list.List
}

func (r *routeNode) routeKey() string {
	return r.key
}

func (r *routeNode) route(matchKeys []string, idx int, routePath []string) (node, error) {
	return nodeListRoute(r.child, matchKeys, idx, routePath)
}

func (r *routeNode) addRoute(id uint, info []string, idx int) error {
	matchKey := info[idx]

	if err := router.CheckRouteKey(matchKey); err != nil {
		return err
	}
	if r.key == "" {
		r.key = matchKey
	}
	if r.childMap == nil {
		r.childMap = make(map[string]*list.Element)
	}
	if r.child == nil {
		r.child = list.New()
	}
	nextKey := info[idx+1]

	if idx == len(info)-2 {
		//	the next node is leaf, push the leaf to end.
		r.child.PushBack(&leafNode{
			id:       id,
			endpoint: info[idx+1],
		})
		r.childMap[nextKey] = r.child.Back()
	} else {
		var (
			nextNode node
			err      error
		)
		c, ok := r.childMap[nextKey]
		if ok {
			nextNode = c.Value.(node)
		} else {
			nextNode = &routeNode{}
			defer func() {
				if err == nil {
					r.child.PushBack(nextNode)
					r.childMap[nextKey] = r.child.Back()
				}
			}()
		}
		if err = nextNode.addRoute(id, info, idx+1); err != nil {
			return err
		}
	}
	return nil
}

func (r *routeNode) deleteRoute(info []string, idx int) (uint, error) {
	return deleteRoute(r.childMap, r.child, info, idx)
}
