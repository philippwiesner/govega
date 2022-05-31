// Package helper
//
// provides basic data structures to build more complex data structures used in the compiler code architecture
//
// hashTable.go implements a basic hash table for storing elements based on a string key lookup in a table.
// When the maximum table size is reached, the hash table will rehash into a bigger table increased by the tableSize
// const.

package helper

import "errors"

var (
	tableSize  uint = 16
	hashFactor uint = 31
)

type bucket struct {
	key   string
	value interface{}
	next  *bucket
}

type hashTable struct {
	table       []*bucket
	cap         uint
	len         uint
	BucketCount uint
}

func NewHashTable() *hashTable {
	return newHashTable(tableSize)
}

func newHashTable(ts uint) *hashTable {
	buckets := make([]*bucket, ts, ts)
	hashTable := hashTable{buckets, ts, 0, 0}
	return &hashTable
}

func (t *hashTable) hash(key string) uint {
	var hash uint = 0
	for _, c := range key {
		hash = ((hashFactor)*hash + uint(c)) % t.cap
	}
	return hash
}

func (t *hashTable) Get(key string) (interface{}, error) {
	hash := t.hash(key)
	entry := t.table[hash]
	if entry == nil {
		return nil, errors.New("element not found")
	}
	if entry.key != "" {
		for {
			if entry.key == key {
				return entry.value, nil
			}
			if entry.next == nil {
				return nil, errors.New("element not found")
			}
			entry = entry.next
		}
	}
	return nil, errors.New("empty key")
}

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

func (t *hashTable) Add(key string, value interface{}) {
	if t.cap == t.len {
		t.increase()
	}
	t.addBucket(key, value)
}
