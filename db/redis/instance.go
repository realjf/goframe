package redis

import "github.com/gogf/gf/container/gmap"

var (
	instances = gmap.NewStrAnyMap(true)
)

func Instance(name ...string) *Redis {
	group := DefaultGroupName
	if len(name) > 0 && name[0] != "" {
		group = name[0]
	}
	v := instances.GetOrSetFuncLock(group, func() interface{} {
		if config, ok := GetConfig(group); ok {
			r := New(config)
			r.group = group
			return r
		}
		return nil
	})
	if v != nil {
		return v.(*Redis)
	}
	return nil
}


