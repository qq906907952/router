package hash_router

import (
	"container/list"
	"fmt"
	"github.com/qq906907952/router/router"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type HashRouter struct {
	// id asc
	raw       []*router.RouteInfo
	keyMap    map[string]*list.Element
	routeNode *list.List
	lock      *sync.RWMutex
}

func (h *HashRouter) Init(infos []*router.RouteInfo) error {
	sort.Slice(infos, func(i, j int) bool {
		return infos[i].Id < infos[j].Id
	})

	h.keyMap = make(map[string]*list.Element)
	h.routeNode = list.New()
	h.lock = &sync.RWMutex{}
	for _, i := range infos {
		i := *i
		is := i.ToSlice()
		if err := h.addRoute(is, i.Id); err != nil {
			return err
		}
		h.raw = append(h.raw, &i)
	}
	return nil
}

func (h *HashRouter) Get(info *router.RouteInfo) (routePath []string, err error) {
	var (
		is   = info.ToSlice()
		leaf node
	)
	routePath = make([]string, 0, len(is))
	h.onReadLockDo(func() {
		if h.routeNode.Len() == 0 {
			err = router.ErrRouteTableNil{}
			return
		}

		leaf, err = nodeListRoute(h.routeNode, is, 0, routePath)
		if err != nil {
			info.Endpoint = ""
		} else {
			leaf := leaf.(*leafNode)
			info.Endpoint = leaf.endpoint
			info.Id = leaf.id
		}

	})
	(*reflect.SliceHeader)(reflect.ValueOf(&routePath).UnsafePointer()).Len = len(is)
	return
}

func (h *HashRouter) Add(info *router.RouteInfo) (err error) {
	i := *info
	is := i.ToSlice()
	for _, v := range is {
		if err = router.CheckRouteKey(v); err != nil {
			return
		}
	}

	// 这里简单起见，不考虑优先级等问题，添加路由只加到最后，并且id自增
	h.onLockDo(func() {
		if len(h.raw) == 0 {
			i.Id = 0
		} else {
			i.Id = h.raw[len(h.raw)-1].Id + 1
		}

		if err := h.addRoute(is, i.Id); err != nil {
			panic("append route fail")
		}
		h.raw = append(h.raw, &i)
	})
	return
}

func (h *HashRouter) Delete(info *router.RouteInfo) (id uint, err error) {
	h.onLockDo(func() {
		if h.routeNode.Len() == 0 {
			err = router.ErrRouteTableNil{}
			return
		}

		delete_is := info.ToSlice()
		id, err = deleteRoute(h.keyMap, h.routeNode, delete_is, 0)
		if err != nil {
			return
		}

		// binary search by id and delete
		start, end := 0, len(h.raw)
		for start < end {
			mid := (start + end) / 2
			if id > h.raw[mid].Id {
				start = mid + 1
			} else {
				end = mid
			}
		}

		for end < len(h.raw) && h.raw[end].Id == id {
			if reflect.DeepEqual(delete_is, h.raw[end].ToSlice()) {
				break
			}
			end++
		}
		h.raw = append(h.raw[:end], h.raw[end+1:]...)
	})
	return
}

func (h *HashRouter) GetRouteTable() string {
	s := make([]string, 0, len(h.raw))
	for _, v := range h.raw {
		s = append(s, strconv.FormatInt(int64(v.Id), 10)+", "+strings.Join(v.ToSlice(), ", "))
	}
	return strings.Join(s, fmt.Sprintln())
}

func (h *HashRouter) onReadLockDo(f func()) {
	h.lock.RLock()
	defer h.lock.RUnlock()
	f()
}

func (h *HashRouter) onLockDo(f func()) {
	h.lock.Lock()
	defer h.lock.Unlock()
	f()
}

func (h *HashRouter) addRoute(is []string, id uint) error {
	var (
		rootNode node
		routeKey = is[0]
	)

	n, fromMap := h.keyMap[routeKey]
	if fromMap {
		rootNode = n.Value.(node)
	} else {
		rootNode = &routeNode{}
	}

	if err := rootNode.addRoute(id, is, 0); err != nil {
		return err
	}
	if !fromMap {
		h.routeNode.PushBack(rootNode)
		h.keyMap[routeKey] = h.routeNode.Back()
	}
	return nil
}
