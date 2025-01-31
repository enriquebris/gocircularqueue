package gocircularqueue

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type Circular struct {
	circularChannel chan string
	mp              sync.Map
	mutex           sync.RWMutex
}

func NewCircularQueue(capacity int) (*Circular, error) {
	ret := &Circular{}
	if err := ret.initialize(capacity); err != nil {
		return nil, err
	}

	return ret, nil
}

func (st *Circular) initialize(capacity int) error {
	if capacity < 1 {
		return fmt.Errorf("capacity must be not less than 1")
	}

	st.circularChannel = make(chan string, capacity)

	return nil
}

func (st *Circular) Length() int {
	return len(st.circularChannel)
}

func (st *Circular) Capacity() int {
	return cap(st.circularChannel)
}

// Enqueue is an atomic enqueue operation
func (st *Circular) Enqueue(key string, value any) (rKey string, rValue any, err error) {
	// locking execution to prevent adding new items at the time Enqueue is being executed
	st.mutex.Lock()
	defer st.mutex.Unlock()

	select {
	case st.circularChannel <- key:
		// queue is not at full capacity

		// store key/value at internal map
		st.mp.Store(key, value)

	default:
		// queue is full

		// dequeues first enqueued pair (no locking operation)
		rKey, rValue, err = st.dequeue(false)
		if err == nil {
			// makes last attempt to enqueue
			select {
			case st.circularChannel <- key:
				// store key/value at internal map
				st.mp.Store(key, value)

			default:
				// unable to enqueue :(
				return rKey, rValue, fmt.Errorf("error while enqueueing key: %v", key)
			}

		} else {
			return rKey, rValue, err
		}
	}

	return rKey, rValue, nil
}

func (st *Circular) dequeue(lock bool) (rKey string, rValue any, err error) {
	if lock {
		// locking execution to prevent any queue update while Dequeue is being executed
		st.mutex.Lock()
		defer st.mutex.Unlock()
	}

	select {
	case rKey = <-st.circularChannel:
		ok := true
		// get key/value from map && delete entry from map
		rValue, ok = st.mp.LoadAndDelete(rKey)
		if !ok {
			// no value found on map for given key
			return rKey, nil, fmt.Errorf("no saved value for key: %v", rKey)
		}

	default:
		return "", nil, errors.New("empty queue")
	}

	return rKey, rValue, nil
}

func (st *Circular) Dequeue() (rKey string, rValue any, err error) {
	return st.dequeue(true)
}

func (st *Circular) Update(key string, value any) error {
	// locking execution to prevent any queue update while Update is being executed
	st.mutex.Lock()
	defer st.mutex.Unlock()

	currentValue, ok := st.mp.Load(key)
	if !ok {
		return fmt.Errorf("key: %v is not currently enqueued", key)
	}

	if !reflect.DeepEqual(currentValue, value) {
		st.mp.Store(key, value)
	}


	return nil
}

func (st *Circular) Get(key string) (rValue any, err error) {
	// locking execution to prevent any queue update while Get is being executed
	st.mutex.RLock()
	defer st.mutex.RUnlock()

	value, ok := st.mp.Load(key)
	if !ok {
		return nil, fmt.Errorf("key: %v is not currently enqueued", key)
	}

	return value, nil
}

func (st *Circular) GetByPrefix(prefix string) (results map[string]any, err error) {
	// locking execution to prevent any queue update while GetByPrefix is being executed
	st.mutex.RLock()
	defer st.mutex.RUnlock()

	results = make(map[string]any)
	st.mp.Range(func(key, value any) bool {
		// check if key is a string and contains the prefix
		if k, ok := key.(string); ok && strings.Contains(k, prefix) {
			results[k] = value
		}
		// continue iteration
		return true
	})

	if len(results) == 0 {
		return nil, fmt.Errorf("no keys found with prefix: %v", prefix)
	}

	return results, nil
}

func (st *Circular) Delete(key string) {
	// locking execution to prevent any queue update while Get is being executed
	st.mutex.RLock()
	defer st.mutex.RUnlock()

	st.mp.Delete(key)
}
