package cache

import (
	"container/list"
	"errors"
	"sync"
)

const (
	THREAD_SAFE_MAP_MAX_CAP = 1000
)

// thread safe map
type IThreadSafeMap interface {
	Add(string, interface{}) error
	Get(string) (interface{}, error)

	CleanAll() error
	Exist(string) bool
	Len() int

	Delete(string) error
	Update(string, interface{}) error
	List() []interface{}
}

type ThreadSafeMap struct {
	lock sync.RWMutex

	items map[string]interface{}

	// 队列
	keys *list.List

	cap int

	len int
}

func NewThreadSafeMap(cap int) IThreadSafeMap {
	if cap <= 0 || cap > THREAD_SAFE_MAP_MAX_CAP {
		cap = THREAD_SAFE_MAP_MAX_CAP
	}
	return &ThreadSafeMap{
		cap:   cap,
		len:   0,
		lock:  sync.RWMutex{},
		items: make(map[string]interface{}, cap),
		keys:  list.New(),
	}
}

func (t *ThreadSafeMap) Add(key string, item interface{}) error {
	if key == "" || item == nil {
		return errors.New("set value error: key or value is empty")
	}

	// 判断长度是否超过限制
	if t.Len() >= THREAD_SAFE_MAP_MAX_CAP {
		lastEle := t.keys.Remove(t.keys.Back()).(string)
		tempMap := map[string]interface{}{}
		// 删除最近最少访问数据
		for k, v := range t.items {
			if k == lastEle {
				t.len--
				continue
			}
			tempMap[k] = v
		}
		t.items = tempMap
	}

	t.lock.Lock()
	defer t.lock.Unlock()

	// 判断是否存在
	t.items[key] = item
	t.len++

	t.keys.PushFront(key)

	return nil
}

func (t *ThreadSafeMap) exist(key string) (bool, *list.Element) {
	for v := t.keys.Front(); v != nil; v = v.Next() {
		if key == v.Value.(string) {
			return true, v
		}
	}

	return false, nil
}

func (t *ThreadSafeMap) Get(key string) (interface{}, error) {
	if key == "" {
		return nil, errors.New("key is empty")
	}

	if t.Len() <= 0 {
		return nil, errors.New("size is empty")
	}

	t.lock.RLock()
	defer t.lock.RUnlock()

	if ele, ok := t.items[key]; ok {
		if ok, item := t.exist(key); ok {
			// 数据访问则把数据移动到头部
			t.keys.MoveToFront(item)
		}

		return ele, nil
	}

	return nil, errors.New(key + " not found")
}

func (t *ThreadSafeMap) List() []interface{} {
	t.lock.RLock()
	defer t.lock.RUnlock()

	data := make([]interface{}, 0, len(t.items))

	for _, v := range t.items {
		data = append(data, v)
	}

	return data
}

// 清理
func (t *ThreadSafeMap) CleanAll() error {
	if t.Len() <= 0 {
		return nil
	}

	t.lock.Lock()
	defer t.lock.Unlock()

	for k, _ := range t.items {
		delete(t.items, k)
		t.len--
	}

	t.keys = list.New()

	return nil
}

func (t *ThreadSafeMap) Exist(key string) bool {
	t.lock.RLock()
	defer t.lock.RUnlock()

	if _, ok := t.items[key]; ok {
		return true
	}

	return false
}

func (t *ThreadSafeMap) Delete(key string) error {
	if key == "" {
		return errors.New("key is empty")
	}
	if t.Len() <= 0 {
		return errors.New("list is empty")
	}
	t.lock.Lock()
	defer t.lock.Unlock()

	delete(t.items, key)
	t.len--

	if ok, data := t.exist(key); ok {
		t.keys.Remove(data)
	}

	return nil
}

func (t *ThreadSafeMap) Update(key string, item interface{}) error {
	if key == "" {
		return errors.New("key is empty")
	}
	t.lock.Lock()
	defer t.lock.Unlock()

	if _, ok := t.items[key]; !ok {
		return errors.New("key[" + key + "] not exist")
	}

	t.items[key] = item

	if ok, item := t.exist(key); ok {
		// 数据访问则把数据移动到头部
		t.keys.MoveToFront(item)
	}

	return nil
}

func (t *ThreadSafeMap) Len() int {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return t.len
}
