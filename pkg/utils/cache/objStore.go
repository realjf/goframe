package cache

import (
	"container/list"
	"errors"
	"sync"
)

const (
	OBJECT_STORE_MAX_CAP = 1000
)

type IObjStore interface {
	Add(interface{}, string, interface{}) error
	Get(interface{}, string) (interface{}, error)

	Exist(interface{}, string) bool
	Len() int

	Delete(interface{}, string) error
	Update(interface{}, string, interface{}) error
	List() []IThreadSafeMap
}

type ObjStore struct {
	cache map[interface{}]IThreadSafeMap
	cap   int
	len   int

	lock sync.RWMutex

	keys *list.List
}

func NewObjStore(cap int) IObjStore {
	if cap <= 0 || cap > OBJECT_STORE_MAX_CAP {
		cap = OBJECT_STORE_MAX_CAP
	}
	return &ObjStore{
		cache: make(map[interface{}]IThreadSafeMap, cap),
		cap:   cap,
		len:   0,
		lock:  sync.RWMutex{},
		keys:  list.New(),
	}
}

func (obj *ObjStore) Add(section interface{}, key string, value interface{}) error {
	if section == nil {
		return errors.New("section is nil")
	}
	if key == "" {
		return errors.New("key is empty")
	}
	if value == nil {
		return errors.New("value is nil")
	}

	// 容量判断
	if obj.Len() >= OBJECT_STORE_MAX_CAP {
		// 删除最近最少访问

	}

	if sectionMap, ok := obj.cache[section]; ok {
		return sectionMap.Add(key, value)
	} else {
		sectionMap := NewThreadSafeMap(THREAD_SAFE_MAP_MAX_CAP)
		err := sectionMap.Add(key, value)
		if err != nil {
			return err
		}
		obj.cache[section] = sectionMap
		obj.len++
	}
	return nil
}

func (obj *ObjStore) Get(section interface{}, key string) (interface{}, error) {
	if section == nil {
		return nil, errors.New("section is nil")
	}
	if key == "" {
		return nil, errors.New("key is empty")
	}

	if data, ok := obj.cache[section]; ok {
		if ok, item := obj.exist(key); ok {
			// 数据访问则把数据移动到头部
			obj.keys.MoveToFront(item)
		}
		return data, nil
	}
	return nil, errors.New(key + " not found")
}

func (obj *ObjStore) Exist(section interface{}, key string) bool {
	if section == nil {
		return false
	}
	if key == "" {
		return false
	}

	_, ok := obj.cache[section]
	return ok
}

func (obj *ObjStore) exist(section string) (bool, *list.Element) {
	for v := obj.keys.Front(); v != nil; v = v.Next() {
		if section == v.Value.(string) {
			return true, v
		}
	}

	return false, nil
}

func (obj *ObjStore) Len() int {
	obj.lock.RLock()
	defer obj.lock.RUnlock()

	return obj.len
}

func (obj *ObjStore) Delete(section interface{}, key string) error {
	if section == nil {
		return errors.New("section is nil")
	}
	if key == "" {
		return errors.New("key is empty")
	}

	if data, ok := obj.cache[section]; !ok {
		return errors.New("section not found")
	} else {
		_, error := data.Get(key)
		if error != nil {
			return error
		} else {
			err := data.Delete(key)
			// 删除节点
			if data.Len() <= 0 {
				delete(obj.cache, section)
				obj.len--
			}
			return err
		}
	}
}

func (obj *ObjStore) Update(section interface{}, key string, value interface{}) error {
	if section == nil {
		return errors.New("section is nil")
	}
	if key == "" {
		return errors.New("key is empty")
	}

	if data, ok := obj.cache[section]; !ok {
		return errors.New("section not found")
	} else {
		_, error := data.Get(key)
		if error != nil {
			return data.Add(key, value)
		} else {
			return data.Update(key, value)
		}
	}
}

func (obj *ObjStore) List() []IThreadSafeMap {
	objList := []IThreadSafeMap{}

	for _, v := range obj.cache {
		objList = append(objList, v)
	}

	return objList
}
