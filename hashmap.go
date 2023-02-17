package byor

import (
	"math"
	"sync"
)

const (
	MAX_LOAD_FACTOR   = 0.75
	MAX_RESIZING_WORK = 128
)

type hnode struct {
	next       *hnode
	key, value string
}

type htable struct {
	table []*hnode
	mask  int32
	size  int32
}

func (ht *htable) keyIndex(key string) int {
	return len(key) & int(ht.mask)
}

func (ht *htable) addEntry(key, val string) {
	newNode := &hnode{
		key:   key,
		value: val,
	}
	index := ht.keyIndex(key)
	node := ht.table[index]
	if node == nil {
		ht.table[index] = newNode
	} else {
		for {
			if node.key == key {
				node.value = val
				break
			} else if node.next == nil {
				node.next = newNode
				ht.size += 1
				break
			} else {
				node = node.next
			}
		}
	}
}

func (ht *htable) getEntry(key string) string {
	index := ht.keyIndex(key)
	node := ht.table[index]
	if node == nil {
		return ""
	} else {
		for {
			if node.key == key {
				return node.value
			} else if node.next == nil {
				return ""
			} else {
				node = node.next
			}
		}
	}
}

func (ht *htable) delEntry(key string) {
	index := ht.keyIndex(key)
	node := ht.table[index]
	if node == nil {
		return
	} else if node.key == key {
		ht.table[index] = node.next
	} else {
		for node.next != nil {
			if node.next.key == key {
				node.next = node.next.next
				break
			} else {
				node = node.next
			}
		}
	}
}

type HMap struct {
	buckets     *htable
	oldbuckets  *htable
	resizingPos int32
	mutex       sync.RWMutex
}

func NewHashMap(size int32) *HMap {
	return &HMap{
		buckets: &htable{
			table: make([]*hnode, size),
			mask:  size - 1,
			size:  0,
		},
		oldbuckets:  nil,
		resizingPos: 0,
		mutex:       sync.RWMutex{},
	}
}

func (hm *HMap) startResize() {
	prevSize := hm.buckets.mask + 1
	if prevSize == math.MaxInt32 {
		return
	}

	var newSize int32
	if prevSize < math.MaxInt32/2 {
		newSize = prevSize * 2
	} else {
		newSize = math.MaxInt32
	}

	hm.buckets, hm.oldbuckets = &htable{
		table: make([]*hnode, newSize),
		mask:  newSize - 1,
		size:  0,
	}, hm.buckets
}

func (hm *HMap) resize() {
	if hm.oldbuckets == nil {
		return
	}

	var maxIndex int32
	for i := hm.resizingPos; i <= maxIndex; i += 1 {
	}

	if hm.oldbuckets.size == 0 {
		hm.oldbuckets = nil
		hm.resizingPos = 0
	}
}

func (hm *HMap) Get(key string) string {
	hm.mutex.RLock()
	defer hm.mutex.RUnlock()

	hm.resize()

	value := hm.buckets.getEntry(key)
	if hm.oldbuckets != nil {
		value = hm.oldbuckets.getEntry(key)
	}
	return value
}

func (hm *HMap) Put(key, val string) {
	hm.mutex.Lock()
	defer hm.mutex.Unlock()
	loadFactor := float32(hm.buckets.size) / float32(hm.buckets.mask+1)
	if loadFactor > MAX_LOAD_FACTOR {
		hm.startResize()
	}
	hm.buckets.addEntry(key, val)
	hm.resize()
}

func (hm *HMap) Delete(key string) {
	hm.mutex.Lock()
	defer hm.mutex.Unlock()

	hm.resize()

	hm.buckets.delEntry(key)
	if hm.oldbuckets != nil {
		hm.oldbuckets.delEntry(key)
	}
}

func (hm *HMap) For(func(key, val string)) {}

func (hm *HMap) Size() int32 {
	hm.mutex.RLock()
	defer hm.mutex.RUnlock()

	size := hm.buckets.size
	if hm.oldbuckets != nil {
		size += hm.oldbuckets.size
	}
	return size
}

func (hm *HMap) Clear() {
	hm.mutex.Lock()
	defer hm.mutex.Unlock()

	hm.buckets = &htable{
		table: make([]*hnode, 2),
		mask:  1,
		size:  0,
	}
	hm.oldbuckets = nil
	hm.resizingPos = 0
}
