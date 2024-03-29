// Package helper
//
// provides basic data structures to build more complex data structures used in the compiler code architecture
//
// hashTable.go implements a basic hash table for storing elements based on a string key lookup in a table.
// When the maximum table size is reached, the hash table will rehash into a bigger table increased by the tableSize
// const.

package helper

type HashTable interface {
	Get(key string) (result interface{}, ok bool)
	Add(key string, value interface{})
	GetBucketCount() uint
	getCapacity() uint
	getLength() uint
}

// tableSize: initial hash table size
//
// hashFactor: prime to calculate hashing value
var (
	tableSize  uint = 128
	hashFactor uint = 31
)

// bucket is an internal structure to store element in hash table.
// on hash collision buckets will be chained together.
type bucket struct {
	key   string      // name of the element in table
	value interface{} // element
	next  *bucket     // link to next bucket
}

// hashTable is a simple table which stored objects via a calculated hash.
//
// Methods for adding new elements are Add and for retrieving Get. If maximum length
// of the table has been reached the table is automatically increased by the tableSize value
// and each bucket is being rehashed.
type hashTable struct {
	table       []*bucket
	cap         uint
	len         uint
	BucketCount uint
}

func (t *hashTable) GetBucketCount() uint {
	return t.BucketCount
}

func (t *hashTable) getCapacity() uint {
	return t.cap
}

func (t *hashTable) getLength() uint {
	return t.len
}

// NewHashTable is the constructor for creating new hashTable pointers. The constructor calls a private
// constructor with the actual tableSize.
func NewHashTable() HashTable {
	var h HashTable = newHashTable(tableSize)
	return h
}

// newHashTable is the private constructor for creating a new hashTable pointer. The private constructor let the
// size of the hashTable choose for better testing.
func newHashTable(ts uint) *hashTable {
	buckets := make([]*bucket, ts, ts)
	return &hashTable{
		table:       buckets,
		cap:         ts,
		len:         0,
		BucketCount: 0,
	}
}

// hash is the private hash function to generate a hash for the keyWord under which the element should be stored.
func (t *hashTable) hash(key string) uint {
	var hash uint = 0
	for _, c := range key {
		hash = ((hashFactor)*hash + uint(c)) % t.cap
	}
	return hash
}

// Get is a function for receiving elements via keyWord hash lookup from the table.
//
// As elements can be of any type an interface is being returned, if the element can not be found
// additional bool status false is being returned.
func (t *hashTable) Get(key string) (result interface{}, ok bool) {
	hash := t.hash(key)
	entry := t.table[hash]
	if entry == nil {
		return nil, false
	}
	if entry.key != "" {
		for {
			if entry.key == key {
				return entry.value, true
			}
			if entry.next == nil {
				return nil, false
			}
			entry = entry.next
		}
	}
	return nil, false
}

// increase is a private method which increases the hashTable via rehashing when the len and cap are equal
// (hashTable is full)
func (t *hashTable) increase() {
	newCap := t.cap + tableSize
	newHashTable := newHashTable(newCap)

	for _, v := range t.table {
		for {
			newHashTable.addBucket(v.key, v.value)
			if v.next == nil {
				break
			}
			v = v.next
		}
	}

	t.table = newHashTable.table
	t.cap = newHashTable.cap
}

// addBucket is a private method to add new elements to the hashTable. This method is used by the increase method
// for rehashing and the public Add method.
func (t *hashTable) addBucket(key string, value interface{}) {
	hash := t.hash(key)
	entry := t.table[hash]
	newBucket := bucket{key, value, nil}
	if entry == nil {
		t.table[hash] = &newBucket
		t.len++
		t.BucketCount++
	} else {
		for {
			if entry.key == key {
				entry.value = value
				break
			}
			if entry.next == nil {
				entry.next = &newBucket
				t.BucketCount++
				break
			}
			entry = entry.next
		}
	}
}

// Add is the public method for adding any value via a key to the hashtable
func (t *hashTable) Add(key string, value interface{}) {
	if t.cap == t.len {
		t.increase()
	}
	t.addBucket(key, value)
}
