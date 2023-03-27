package router

import "fmt"

type ErrInitFail struct {
	reason string
}

func (i ErrInitFail) Error() string {
	return i.reason
}

func (i ErrInitFail) New(f string, v ...interface{}) ErrInitFail {
	return ErrInitFail{
		reason: fmt.Sprintf(f, v...),
	}
}

type ErrUnReachable struct {
}

func (e ErrUnReachable) Error() string {
	return fmt.Sprintf("route unreachable")
}

type ErrInfoNullKey struct {
}

func (e ErrInfoNullKey) Error() string {
	return "route key can not null"
}

type ErrInfoEndpointFormat struct {
}

func (e ErrInfoEndpointFormat) Error() string {
	return "endpoint must nil or start with * or end with *"
}

type ErrRouteKeyNotExist struct {
	key string
}

func (e ErrRouteKeyNotExist) New(key string) ErrRouteKeyNotExist {
	return ErrRouteKeyNotExist{
		key: key,
	}
}
func (e ErrRouteKeyNotExist) Error() string {
	return fmt.Sprintf("route key %s  not found", e.key)
}

type ErrRouteTableNil struct {
}

func (e ErrRouteTableNil) Error() string {
	return fmt.Sprintf("table is nil")
}
