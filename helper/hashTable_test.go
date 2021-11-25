package helper

import (
	"reflect"
	"testing"
)

func TestHashTable(t *testing.T) {
	tests := []struct {
		key   string
		value interface{}
	}{
		{"hello", 1234},
		{"world", "fsfgdfsdfgsfsd"},
	}

	hashTable := New()
	cap := cap(hashTable.table)

	if cap != int(tableSize) {
		t.Fatalf("Cap is %v, should be %v", cap, tableSize)
	}

	for _, tc := range tests {
		hashTable.Add(tc.key, tc.value)
	}

	for i, tc := range tests {
		got := hashTable.Get(tc.key)
		if !reflect.DeepEqual(tc.value, got) {
			t.Fatalf("test %d: expected: %v, got: %v for key: %v! BucketCount: %d", i+1, tc.value, got, tc.key, hashTable.BucketCount)
		}
	}

}
