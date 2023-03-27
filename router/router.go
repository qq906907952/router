package router

import "strings"

type RouteInfo struct {
	Id       uint
	Domain   string
	Service  string
	Method   string
	Endpoint string
}

func (r *RouteInfo) ToSlice() []string {
	return []string{
		strings.TrimSpace(r.Domain),
		strings.TrimSpace(r.Service),
		strings.TrimSpace(r.Method),
		strings.TrimSpace(r.Endpoint),
	}
}

type Router interface {
	Init([]*RouteInfo) error
	Get(*RouteInfo) (routePath []string, err error) // Fill the RouteInfo.Endpoint and RouteInfo.Id by the route info
	Add(info *RouteInfo) (err error)
	Delete(info *RouteInfo) (id uint, err error)
}

func CheckRouteKey(k string) error {
	if k == "" {
		return ErrInfoNullKey{}
	}
	return nil
}
