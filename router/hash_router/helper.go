package hash_router

import (
	"container/list"
	"github.com/qq906907952/router/router"
	"regexp"
	"strings"
)

/*
return whether matchKey match routeKey.
routeKey start with * will act as suffix match, end with * will act as prefix match, * will match any, otherwise will act as regexp.
*/
func matchRouteKey(matchKey, routeKey string) (bool, error) {
	var (
		err   error
		match = func() bool {
			if routeKey == "*" {
				return true
			} else if strings.HasSuffix(routeKey, "*") {
				// prefix match, such as xxx*

				return strings.HasPrefix(matchKey, routeKey[:len(routeKey)-1])

			} else if strings.HasPrefix(matchKey, "*") {

				// suffix match, such as *xxx
				return strings.HasSuffix(matchKey, routeKey[1:])

			} else {
				if matchKey == routeKey {
					return true
				}
				// 这里可以优化成regexp预编译，这里纯属偷懒
				var reg *regexp.Regexp
				reg, err = regexp.Compile("\\b" + routeKey + "\\b")
				if err != nil {
					return false
				}
				return reg.MatchString(matchKey)
			}
		}
	)

	if err := router.CheckRouteKey(matchKey); err != nil {
		return false, err
	}
	return match(), err
}

// recursive to match route
func nodeListRoute(nodeList *list.List, matchKeys []string, idx int, routePath []string) (node, error) {
	matchKey := matchKeys[idx]
	if matchKey == "" {
		return nil, router.ErrInfoNullKey{}
	}

	var (
		matchNode node
		err       error
		match     bool
	)

	for e := nodeList.Front(); e != nil; e = e.Next() {
		childNode := e.Value.(node)

		match, err = matchRouteKey(matchKey, childNode.routeKey())
		if err != nil {
			return nil, err
		}
		if match {
			// routePath append is safe because slice is pre-alloc memory
			routePath = append(routePath, childNode.routeKey())

			// maybe have same route path to different endpoint, return the first endpoint when recursive to the latest route node
			if idx == len(matchKeys)-2 {
				leaf := childNode.(*routeNode).child.Front().Value.(node)
				routePath = append(routePath, leaf.routeKey())
				return leaf, nil
			}
			matchNode, err = childNode.route(matchKeys, idx+1, routePath)

			if err != nil {
				// routePath backtrace
				routePath = routePath[:len(routePath)-1]
				// the current route path not match, continue next route key
				if _, ok := err.(router.ErrUnReachable); ok {
					continue
				}
				return nil, err
			}
			return matchNode, nil
		} else {
			continue
		}
	}
	return nil, router.ErrUnReachable{}
}

// backtrace delete
func deleteRoute(childMap map[string]*list.Element, childList *list.List, info []string, idx int) (uint, error) {
	matchKey := info[idx]
	child, ok := childMap[matchKey]
	if !ok {
		return 0, router.ErrRouteKeyNotExist{}.New(matchKey)
	}
	id, err := child.Value.(node).deleteRoute(info, idx+1)
	if err != nil {
		return 0, err
	}
	childList.Remove(child)
	delete(childMap, matchKey)
	return id, nil
}
