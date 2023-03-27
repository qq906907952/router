package hash_router

import (
	"fmt"
	"github.com/qq906907952/router/router"
	"reflect"
	"testing"
)

var rt = []*router.RouteInfo{

	{
		Id:       4,
		Domain:   "aa*",
		Service:  "*bb",
		Method:   "cc*",
		Endpoint: "dd*",
	},

	{
		Id:       0,
		Domain:   "a",
		Service:  "b",
		Method:   "c",
		Endpoint: "d",
	},

	{
		Id:       1,
		Domain:   "a",
		Service:  "bb*",
		Method:   "*cc",
		Endpoint: "ddbc",
	},

	{
		Id:       1,
		Domain:   "aa",
		Service:  "bb",
		Method:   "cc",
		Endpoint: "dd",
	},

	{
		Id:       2,
		Domain:   "aa",
		Service:  "bb",
		Method:   "cc",
		Endpoint: "dd",
	},
}

const table = `0, a, b, c, d
1, a, bb*, *cc, ddbc
1, aa, bb, cc, dd
2, aa, bb, cc, dd
4, aa*, *bb, cc*, dd*`

func must(f func() (bool, string)) {
	s, b := f()
	if !s {
		panic(b)
	}
}

func TestHashRouter(t *testing.T) {
	r := HashRouter{}

	must(func() (bool, string) {
		err := r.Init(rt)
		return err == nil, fmt.Sprintf("init fail %v", err)
	})

	must(func() (bool, string) {
		return r.GetRouteTable() == table,
			fmt.Sprintf("GetRouteTable test fail, expect: %s%s%sactually:%s%s", fmt.Sprintln(), table, fmt.Sprintln(), fmt.Sprintln(), r.GetRouteTable())
	})

	must(func() (bool, string) {
		i := &router.RouteInfo{
			Domain:  "a",
			Service: "bbaaaaaaaa",
			Method:  "asdasdsacc",
		}
		p, err := r.Get(i)

		must(func() (bool, string) {
			return err == nil, fmt.Sprintf("get fail: %v", err)
		})

		must(func() (bool, string) {
			return i.Endpoint == "ddbc", fmt.Sprintf("endpoint expect %s, actually %s", "ddbc", i.Endpoint)
		})

		return reflect.DeepEqual(p, rt[1].ToSlice()), fmt.Sprintf("route path expect %v actually %v", rt[1].ToSlice(), p)
	})

	addAndDelete := &router.RouteInfo{
		Domain:   "z",
		Service:  "zz*",
		Method:   "*zz",
		Endpoint: "ffff",
	}
	// add
	must(func() (bool, string) {

		must(func() (bool, string) {
			err := r.Add(addAndDelete)
			return err == nil, fmt.Sprintf("add fail: %v", err)
		})
		i := &router.RouteInfo{
			Domain:  "z",
			Service: "zzdddddddddddddd",
			Method:  "ddddddddddddddddddddddddzz",
		}
		p, err := r.Get(i)

		must(func() (bool, string) {
			return err == nil, fmt.Sprintf("get fail after delete: %v", err)
		})

		addAndDelete.Id = i.Id

		must(func() (bool, string) {
			return i.Endpoint == "ffff", fmt.Sprintf("endpoint expect %s, actually %s", "ffff", i.Endpoint)
		})

		return reflect.DeepEqual(p, addAndDelete.ToSlice()), fmt.Sprintf("route path expect %v actually %v", addAndDelete.ToSlice(), p)
	})

	// delete
	must(func() (bool, string) {
		id, err := r.Delete(addAndDelete)
		must(func() (bool, string) {
			return err == nil, fmt.Sprintf("delete fail: %v", err)
		})

		must(func() (bool, string) {
			return id == addAndDelete.Id, "delete return unexpect id"
		})

		must(func() (bool, string) {
			return r.GetRouteTable() == table,
				fmt.Sprintf("GetRouteTable test fail, expect: %s%s%sactually:%s%s", fmt.Sprintln(), table, fmt.Sprintln(), fmt.Sprintln(), r.GetRouteTable())
		})

		_, err = r.Get(addAndDelete)
		_, ok := err.(router.ErrUnReachable)
		return ok, "almost impossible"
	})
}
