package gocircularqueue

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
// ** Delete operations
// ***************************************************************************************
// TestDeleteExistingKey delete an existing key
func (suite *CircularTestSuite) TestDeleteExistingKey() {
	// enqueue single pair
	suite.circular.Enqueue("0", 0)

	// delete key
	suite.circular.Delete("0")

	// Verify that the key no longer exists
	rValue, err := suite.circular.Get("0")
	// Error is expected since the key should not exist
	suite.Error(err, "expected an error when getting a deleted key")
	suite.Nil(rValue, "expected the return value to be nil for a deleted key")
}

// TestDeleteNonExistingKey tests that deleting a non-existing key does not cause any errors or side effects.
func (suite *CircularTestSuite) TestDeleteNonExistingKey() {
	// Arrange: ensure the cache is empty
	suite.Equal(0, suite.circular.Length(), "expected cache to be empty before test")

	// Act: try to delete a non-existing key
	suite.circular.Delete("nonexistent")

	// Assert: Verify the cache is still empty
	suite.Equal(0, suite.circular.Length(), "expected cache to remain empty after deleting a non-existing key")
}

// ***************************************************************************************
// ** GetByPrefix operations
// ***************************************************************************************
// TestGetByPrefixEmptyCache tests that GetByPrefix returns an empty map when the cache is empty.
func (suite *CircularTestSuite) TestGetByPrefixEmptyCache() {
	// Act: get all key-value pairs with a prefix
	pairs, err := suite.circular.GetByPrefix("prefix")

	// Assert: Verify the result is an empty map
	suite.Error(err, "expected an error when getting by prefix in an empty cache")
	suite.Nil(pairs, "expected an Nil when getting by prefix in an empty cache")
}

// TestGetByPrefixNoMatchingKeys tests that GetByPrefix returns an empty map when there are no keys with the given prefix.
func (suite *CircularTestSuite) TestGetByPrefixNoMatchingKeys() {
	// Arrange: enqueue a single key-value pair
	suite.circular.Enqueue("other-prefix", "prefix-value")

	// Act: get all key-value pairs with a prefix
	pairs, err := suite.circular.GetByPrefix("u.1-a.2")

	// Assert: Verify the result is an error and the map is nil
	suite.Error(err, "expected an error when getting by prefix with no matching keys")
	suite.Nil(pairs, "expected a nil map when getting by prefix with no matching keys")
}

// TestGetByPrefixSingleMatchingKey tests that GetByPrefix returns a map with a single key-value pair when there is a single key with the given prefix.
func (suite *CircularTestSuite) TestGetByPrefixSingleMatchingKey() {
	// Arrange: enqueue a single key-value pair
	suite.circular.Enqueue("prefix-key", "prefix-value")
	suite.circular.Enqueue("other-key", "other-value")

	// Act: get all key-value pairs with a prefix
	pairs, err := suite.circular.GetByPrefix("prefix")

	// Assert: Verify the result contains the single key-value pair
	suite.NoError(err, "expected no error when getting by prefix in a cache with a single matching key")
	suite.Len(pairs, 1, "expected a single key-value pair when getting by prefix in a cache with a single matching key")
	suite.Contains(pairs, "prefix-key", "expected the key to be present in the result")
	suite.Equal("prefix-value", pairs["prefix-key"], "expected the value to be correct in the result")
}

// TestGetByPrefixMultipleMatchingKeys tests that GetByPrefix returns a map with all key-value pairs that have the given prefix.
func (suite *CircularTestSuite) TestGetByPrefixMultipleMatchingKeys() {
	// Arrange: enqueue multiple key-value pairs
	suite.circular.Enqueue("prefix-key-1", "prefix-value-1")
	suite.circular.Enqueue("prefix-key-2", "prefix-value-2")
	suite.circular.Enqueue("other-key", "other-value")

	// Act: get all key-value pairs with a prefix
	pairs, err := suite.circular.GetByPrefix("prefix")

	// Assert: Verify the result contains all matching key-value pairs
	suite.NoError(err, "expected no error when getting by prefix in a cache with multiple matching keys")
	suite.Len(pairs, 2, "expected all matching key-value pairs when getting by prefix in a cache with multiple matching keys")
	suite.Contains(pairs, "prefix-key-1", "expected the first key to be present in the result")
	suite.Equal("prefix-value-1", pairs["prefix-key-1"], "expected the first value to be correct in the result")
	suite.Contains(pairs, "prefix-key-2", "expected the second key to be present in the result")
	suite.Equal("prefix-value-2", pairs["prefix-key-2"], "expected the second value to be correct in the result")
}

// TestGetByPrefixEmptyPrefix tests getting keys with an empty prefix, which should return all keys in the cache.
func (suite *CircularTestSuite) TestGetByPrefixEmptyPrefix() {
	// enqueue pairs with different prefixes
	suite.circular.Enqueue("u.1-a.2", "value1")
	suite.circular.Enqueue("u.2-a.2", "value2")
	suite.circular.Enqueue("u.3-a.3", "value3")

	// get keys with empty prefix
	results, err := suite.circular.GetByPrefix("")
	// no error expected
	suite.NoError(err)
	// expected results (all keys)
	expected := map[string]any{
		"u.1-a.2": "value1",
		"u.2-a.2": "value2",
		"u.3-a.3": "value3",
	}
	suite.Equal(expected, results)
}

// ***************************************************************************************
// ** Run suite
// ***************************************************************************************

func TestCircularTestSuite(t *testing.T) {
	suite.Run(t, new(CircularTestSuite))
}
