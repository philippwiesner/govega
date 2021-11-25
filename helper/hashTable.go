package helper

const (
	tableSize  uint8 = 4
	hashFactor uint8 = 31
)

type bucket struct {
	key   string
	value interface{}
	next  *bucket
}

type HashTable struct {
	table       []bucket
	cap         uint
	len         uint
	BucketCount uint
}

func New() *HashTable {
	buckets := make([]bucket, tableSize, tableSize)
	hashTable := HashTable{buckets, uint(tableSize), 0, 0}
	return &hashTable
}

func (t *HashTable) hash(key string) uint {
	var hash uint = 0
	for _, c := range key {
		hash = (uint(hashFactor)*hash + uint(c)) % t.cap
	}
	return hash
}

func (t *HashTable) increase() {
	newCap := t.cap + uint(tableSize)
	newTable := make([]bucket, newCap, newCap)

	for _, v := range t.table {
		hash := t.hash(v.key)
		newTable[hash] = v
	}

	t.table = newTable
	t.cap = newCap
}

func (t *HashTable) Get(key string) interface{} {
	hash := t.hash(key)
	entry := t.table[hash]
	if entry.key != "" {
		for {
			if entry.key == key {
				return entry.value
			}
			if entry.next == nil {
				return nil
			}
			entry = *entry.next
		}
	}
	return nil
}

func (t *HashTable) Add(key string, value interface{}) {
	hash := t.hash(key)
	entry := t.table[hash]
	newBucket := bucket{key, value, nil}
	if entry.key == "" {
		if t.cap == t.len {
			t.increase()
		}
		t.table[hash] = newBucket
		t.len++
		t.BucketCount++
	} else {
		for {
			if entry.key == key {
				t.table[hash].value = value
				break
			}
			if entry.next == nil {
				t.table[hash].next = &newBucket
				t.BucketCount++
				break
			}
			entry = *entry.next
		}
	}
}
