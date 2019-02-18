package cache

import "testing"

func TestNewObjStore(t *testing.T) {
	var objStore IObjStore
	objStore = NewObjStore(1000)
	err := objStore.Add("name", "12313", "1233")
	if err != nil {
		t.Fatal(err)
	}
	data, err := objStore.Get("name", "12313")
	if err != nil {
		t.Fatal(err)
	}
	t.Fatal(data)
}



