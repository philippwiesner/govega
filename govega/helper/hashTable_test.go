package helper

import (
	"math/rand"
	"testing"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type testBucket struct {
	key   string
	value string
}

type testWant struct {
	BucketCount uint
	TableSize   uint
}

func randomString(l int) string {
	bytes := make([]rune, l)

	for i := range bytes {
		bytes[i] = letters[rand.Intn(len(letters))]
	}

	return string(bytes)
}

func randomTestData(len int) []testBucket {
	testBuckets := make([]testBucket, len)

	min := 5
	max := 20

	for i := 0; i < len; i++ {
		testBuckets[i].key = randomString(rand.Intn(max-min+1) + min)
		testBuckets[i].value = randomString(rand.Intn(max-min+1) + min)
	}

	return testBuckets
}

func TestHashTable_Add(t *testing.T) {
	tests := []struct {
		in   []testBucket
		want testWant
	}{
		{randomTestData(10), testWant{10, 128}},
		{randomTestData(200), testWant{200, 128}},
	}

	for i, tc := range tests {

		hashTable := NewHashTable()

		for _, b := range tc.in {
			hashTable.Add(b.key, b.value)
		}

		gotBuckets := hashTable.BucketCount
		if gotBuckets != tc.want.BucketCount {
			t.Fatalf("test %d: Capacitiy expected: %v, got: %v", i+1, tc.want.BucketCount, gotBuckets)
		}

		gotTableSize := hashTable.cap
		if gotTableSize != tc.want.TableSize {
			t.Fatalf("test %d: TableSize expected: %v, got :%v. Debug: Table length %v", i+1, tc.want.TableSize, gotTableSize, hashTable.len)
		}

		for _, b := range tc.in {
			got, _ := hashTable.Get(b.key)
			if got != b.value {
				t.Fatalf("test %d: expected: %v:%v, got %v:%v", i+1, b.key, b.value, b.key, got)
			}
		}

	}
}

func TestHashTable_Get(t *testing.T) {
	ht := NewHashTable()
	el, _ := ht.Get("")
	if el != nil {
		t.Fatalf("Empty HashTable should return nil element")
	}
}
