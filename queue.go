package gocircularqueue

type CircularQueue interface {
	// Capacity returns queue capacity
	Capacity() int
	// Length return total amount of enqueued items
	Length() int
	// Enqueue enqueues a new item into the queue
	Enqueue(key string, value any) (rKey string, rValue any, err error)
	// Dequeue dequeues the first enqueued key/value pair
	Dequeue() (rKey string, rValue any, err error)
	// Update updates a key/value pair based on the given key
	Update(key string, value any) error
	// Get returns the key/value for a pair currently enqueued
	Get(key string) (rValue any, err error)
}
