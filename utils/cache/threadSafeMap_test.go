package cache

import (
	"testing"
	"strconv"
)

func TestNewThreadSafeMap(t *testing.T) {
	var threadMap IThreadSafeMap

	threadMap = NewThreadSafeMap(100)

	for i:= 1; i <= 1000; i++ {
		threadMap.Add(strconv.Itoa(i), i)
	}

	threadMap.Update("2", "10001")

	_, err := threadMap.Get("1")
	if err != nil {
		t.Fatal(err.Error())
	}else{
		threadMap.Add("1001", 1001)
	}

	data, _ := threadMap.Get("3")

	threadMap.Add("1002", 1002)
	data, _ = threadMap.Get("4")

	t.Fatal(threadMap.List(), threadMap.Len(), data)
}


