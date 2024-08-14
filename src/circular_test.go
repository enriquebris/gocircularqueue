package src

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"math/rand/v2"
	"testing"
)

const (
	capacity = 100
)

type CircularTestSuite struct {
	suite.Suite
	circular *Circular
}

func (suite *CircularTestSuite) SetupTest() {
	var err error
	suite.circular, err = NewCircularQueue(capacity)
	suite.NoError(err)
}

// ***************************************************************************************
// ** Queue initialization
// ***************************************************************************************

// TestNoEnqueuedItemsAtInitialization verifies no items are enqueued on a fresh circular queue
func (suite *CircularTestSuite) TestNoEnqueuedItemsAtInitialization() {
	suite.Equalf(0, suite.circular.Length(), "")
}

// TestZeroCapacityNewQueue makes an attempt to generate a zero capacity circular queue
func (suite *CircularTestSuite) TestZeroCapacityNewQueue() {
	cQueue, err := NewCircularQueue(0)
	// error expected
	suite.Error(err)
	// nil queue expected
	suite.Nil(cQueue)
}

// ***************************************************************************************
// ** Enqueue / Dequeue operations
// ***************************************************************************************

// TestEnqueueSingleItemOnEmptyQueue enqueues a single item on an empty queue, then
// dequeues it
func (suite *CircularTestSuite) TestEnqueueSingleItemOnEmptyQueue() {
	const (
		key   = "1"
		value = 1
	)
	// enqueue a single item
	rKey, rValue, err := suite.circular.Enqueue(key, value)
	// no error expected
	suite.NoError(err)
	// no returned key-value pair expected
	suite.Empty(rKey)
	suite.Nil(rValue)

	// length == 1
	suite.Equal(1, suite.circular.Length())

	// capacity == capacity
	suite.Equal(capacity, suite.circular.Capacity())

	// dequeue enqueued item
	k, v, err := suite.circular.Dequeue()
	// no error expected
	suite.NoError(err)
	// validating enqueued key/value pair
	suite.Equal(key, k)
	suite.Equal(value, v)
}

// TestEnqueueFullQueue enqueues all possible items (capacity) on an empty queue, then enqueues an extra item
func (suite *CircularTestSuite) TestEnqueueFullQueue() {
	for i := 0; i < capacity; i++ {
		keyToEnqueue := fmt.Sprintf("%v", i)
		rKey, rValue, err := suite.circular.Enqueue(keyToEnqueue, i)
		// no error expected
		suite.NoError(err)
		// no returned key-value pair expected
		suite.Empty(rKey)
		suite.Nil(rValue)

		// capacity == capacity
		suite.Equal(capacity, suite.circular.Capacity())
		// length
		suite.Equal(i+1, suite.circular.Length())

		// get enqueued item
		rValue2, err := suite.circular.Get(keyToEnqueue)
		// no error expected
		suite.NoError(err)
		// value should be equal to enqueued value (i)
		suite.Equal(i, rValue2)
	}

	// enqueueing extra item (capacity + 1)
	rKey, rValue, err := suite.circular.Enqueue("extra", "extra")
	// no error expected
	suite.NoError(err)
	// returned key-pair expected (first pair enqueued: 0:0)
	suite.Equal("0", rKey)
	suite.Equal(0, rValue)

	// verify rKey is not anymore in the queue
	rValue2, err := suite.circular.Get(rKey)
	// error expected
	suite.Error(err)
	// no value expected
	suite.Nil(rValue2)
}

// TestDequeueOnEmptyQueue dequeues on an empty queue
func (suite *CircularTestSuite) TestDequeueOnEmptyQueue() {
	rKey, rValue, err := suite.circular.Dequeue()
	// error expected
	suite.Error(err)
	// empty key/value expected
	suite.Empty(rKey)
	suite.Nil(rValue)
}

// ***************************************************************************************
// ** Update operations
// ***************************************************************************************

// TestUpdateOnEmptyQueue makes an attempt to update a key/value pair on an empty queue
func (suite *CircularTestSuite) TestUpdateOnEmptyQueue() {
	err := suite.circular.Update("0", 0)
	// error expected
	suite.Error(err)
}

// TestUpdateMissingItemOnFullQueue makes an attempt to update a missing key/value pair on
// a full queue
func (suite *CircularTestSuite) TestUpdateMissingItemOnFullQueue() {
	// enqueueing max items on queue
	for i := 0; i < capacity; i++ {
		suite.circular.Enqueue(fmt.Sprintf("%v", i), i)
	}

	err := suite.circular.Update("missing-key", "any-value")
	// error expected
	suite.Error(err)
}

// TestUpdateSingleItem updates a single -enqueued- pair
func (suite *CircularTestSuite) TestUpdateSingleItem() {
	// enqueue single pair
	suite.circular.Enqueue("0", 0)

	// update pair on queue
	err := suite.circular.Update("0", "updated-value")
	// no error expected
	suite.NoError(err)

	// verify updated pair
	rValue, err := suite.circular.Get("0")
	// no error expected
	suite.NoError(err)
	suite.Equal("updated-value", rValue)
}

// TestUpdateRandomPair updates a random enqueued pair
func (suite *CircularTestSuite) TestUpdateRandomPair() {
	// enqueue max amount of pairs
	for i := 0; i < capacity; i++ {
		suite.circular.Enqueue(fmt.Sprintf("%v", i), i)
	}

	// update random pair
	index := rand.IntN(capacity - 1)
	key := fmt.Sprintf("%v", index)
	// update
	err := suite.circular.Update(key, "updated-value")
	// no error expected
	suite.NoError(err)

	// verify updated pair
	rValue, err := suite.circular.Get(key)
	// no error expected
	suite.NoError(err)
	suite.Equal("updated-value", rValue)
}

// ***************************************************************************************
// ** Run suite
// ***************************************************************************************

func TestCircularTestSuite(t *testing.T) {
	suite.Run(t, new(CircularTestSuite))
}
