package src

import (
	"fmt"
	"github.com/stretchr/testify/suite"
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

func (suite *CircularTestSuite) TestNoEnqueuedItemsAtInitialization() {
	suite.Equalf(0, suite.circular.Length(), "")
}

// ***************************************************************************************
// ** Enqueue operations
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
	suite.Equal("", rKey)
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
		rKey, rValue, err := suite.circular.Enqueue(fmt.Sprintf("%v", i), i)
		// no error expected
		suite.NoError(err)
		// no returned key-value pair expected
		suite.Equal("", rKey)
		suite.Nil(rValue)

		// capacity == capacity
		suite.Equal(capacity, suite.circular.Capacity())
		// length
		suite.Equal(i+1, suite.circular.Length())
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

// ***************************************************************************************
// ** Run suite
// ***************************************************************************************

func TestCircularTestSuite(t *testing.T) {
	suite.Run(t, new(CircularTestSuite))
}
