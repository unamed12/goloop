package trie_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/icon-project/goloop/common/db"
	"github.com/icon-project/goloop/common/trie"
	"github.com/icon-project/goloop/common/trie/mpt"
)

type testDB struct {
	pool map[string][]byte
}

func (db *testDB) Get(k []byte) ([]byte, error) {
	return db.pool[string(k)], nil
}

func (db *testDB) Set(k, v []byte) error {
	db.pool[string(k)] = v
	return nil
}

func (db *testDB) Batch() db.Batch {

	return nil
}
func (db *testDB) Has(key []byte) bool {
	return false
}

func (db *testDB) Delete(key []byte) error {

	return nil
}

func (db *testDB) Transaction() (db.Transaction, error) {
	return nil, nil
}

func (db *testDB) Iterator() db.Iterator {
	return nil
}

func (db *testDB) Close() error {
	return nil
}

func newDB() *testDB {
	return &testDB{pool: make(map[string][]byte)}

}

var testPool = map[string]string{
	"doe":          "reindeer",
	"dog":          "puppy",
	"dogglesworth": "cat",
}

func TestInsert(t *testing.T) {
	db := newDB()
	manager := mpt.NewManager(db)
	trie := manager.NewMutable(nil)

	for k, v := range testPool {
		updateString(trie, k, v)
	}

	hashHex := "8aad789dff2f538bca5d8ea56e8abe10f4c7ba3a5dea95fea4cd6e7c3a1168d3"
	strRoot := fmt.Sprintf("%x", trie.RootHash())
	if strings.Compare(strRoot, hashHex) != 0 {
		t.Errorf("exp %s got %s", hashHex, strRoot)
	}

	immutable := trie.GetSnapshot()
	immutable.Flush()

	mutable := manager.NewMutable(nil)
	mutable.Reset(immutable)
	doeV, _ := mutable.Get([]byte("doe"))
	if strings.Compare(testPool["doe"], string(doeV)) != 0 {
		t.Errorf("%s vs %s", testPool["doe"], string(doeV))
	}

	trie = manager.NewMutable(nil)
	updateString(trie, "A", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")

	hashHex = "d23786fb4a010da3ce639d66d5e904a11dbc02746d1ce25029e53290cabf28ab"
	strRoot = fmt.Sprintf("%x", trie.RootHash())
	if strings.Compare(strRoot, hashHex) != 0 {
		t.Errorf("exp %s got %s", hashHex, strRoot)
	}
}

func TestDelete1(t *testing.T) {
	db := newDB()
	manager := mpt.NewManager(db)
	trie := manager.NewMutable(nil)

	updateString(trie, "doe", "reindeer")
	solution1 := fmt.Sprintf("%x", trie.RootHash())
	updateString(trie, "dog", "puppy")
	solution2 := fmt.Sprintf("%x", trie.RootHash())
	updateString(trie, "dogglesworth", "cat")

	hashHex := "8aad789dff2f538bca5d8ea56e8abe10f4c7ba3a5dea95fea4cd6e7c3a1168d3"
	strRoot := fmt.Sprintf("%x", trie.RootHash())
	if strings.Compare(strRoot, hashHex) != 0 {
		t.Errorf("exp %s got %s", hashHex, strRoot)
	}

	trie.Delete([]byte("dogglesworth"))
	resultRoot := fmt.Sprintf("%x", trie.RootHash())
	if strings.Compare(solution2, resultRoot) != 0 {
		t.Errorf("solution %s, result %s", solution2, resultRoot)
	}
	trie.Delete([]byte("dog"))
	resultRoot = fmt.Sprintf("%x", trie.RootHash())
	if strings.Compare(solution1, resultRoot) != 0 {
		t.Errorf("solution %s, result %s", solution1, resultRoot)
	}
}

func TestDelete2(t *testing.T) {
	db := newDB()
	manager := mpt.NewManager(db)
	trie := manager.NewMutable(nil)
	vals := []struct{ k, v string }{
		{"do", "verb"},
		{"ether", "wookiedoo"},
		{"horse", "stallion"},
		{"shaman", "horse"},
		{"doge", "coin"},
		{"ether", ""},
		{"dog", "puppy"},
		{"shaman", ""},
	}
	for _, val := range vals {
		if val.v != "" {
			updateString(trie, val.k, val.v)
		} else {
			deleteString(trie, val.k)
		}
	}

	strRoot := fmt.Sprintf("%x", trie.RootHash())
	hashHex := "5991bb8c6514148a29db676a14ac506cd2cd5775ace63c30a4fe457715e9ac84"
	if strings.Compare(strRoot, hashHex) != 0 {
		t.Errorf("exp %s got %s", hashHex, strRoot)
	}
}

func TestCache(t *testing.T) {
	db := newDB()
	manager := mpt.NewManager(db)
	mutable := manager.NewMutable(nil)

	for k, v := range testPool {
		updateString(mutable, k, v)
	}

	hashHex := "8aad789dff2f538bca5d8ea56e8abe10f4c7ba3a5dea95fea4cd6e7c3a1168d3"
	root := mutable.RootHash()
	strRoot := fmt.Sprintf("%x", root)
	if strings.Compare(strRoot, hashHex) != 0 {
		t.Errorf("exp %s got %s", hashHex, strRoot)
	}

	snapshot := mutable.GetSnapshot()
	snapshot.Flush()

	// check : Does db in Snapshot have to be passed to Mutable?
	//cacheTrie := mpt.NewCache(nil)
	//cacheTrie.Load(db, root)
	immutable := manager.NewImmutable(root)
	for k, v := range testPool {
		value, _ := immutable.Get([]byte(k))
		if bytes.Compare(value, []byte(v)) != 0 {
			t.Errorf("Wrong value. expected [%x] but [%x]", v, value)
		}
	}

}

func TestDeleteSnapshot(t *testing.T) {
	// delete, snapshot, write
	db := newDB()
	manager := mpt.NewManager(db)
	trie := manager.NewMutable(nil)

	updateString(trie, "doe", "reindeer")
	updateString(trie, "dog", "puppy")
	solution2 := fmt.Sprintf("%x", trie.RootHash())
	updateString(trie, "dogglesworth", "cat")

	hashHex := "8aad789dff2f538bca5d8ea56e8abe10f4c7ba3a5dea95fea4cd6e7c3a1168d3"
	strRoot := fmt.Sprintf("%x", trie.RootHash())
	if strings.Compare(strRoot, hashHex) != 0 {
		t.Errorf("exp %s got %s", hashHex, strRoot)
	}

	trie.Delete([]byte("dogglesworth"))
	resultRoot := fmt.Sprintf("%x", trie.RootHash())
	if strings.Compare(solution2, resultRoot) != 0 {
		t.Errorf("solution %s, result %s", solution2, resultRoot)
	}

	snapshot := trie.GetSnapshot()
	snapshot.Flush()
	trie2 := manager.NewMutable(nil)
	trie2.Reset(snapshot)
	updateString(trie, "dogglesworth", "cat")
	updateString(trie2, "dogglesworth", "cat")
	strRoot = fmt.Sprintf("%x", trie2.RootHash())
	if strings.Compare(strRoot, hashHex) != 0 {
		t.Errorf("exp %s got %s", hashHex, strRoot)
	}
	trie2.Delete([]byte("dogglesworth"))
	solution3 := fmt.Sprintf("%x", trie2.RootHash())
	if strings.Compare(solution3, resultRoot) != 0 {
		t.Errorf("solution %s, result %s", solution3, resultRoot)
	}
	updateString(trie2, "dogglesworth", "cat")
	strRoot = fmt.Sprintf("%x", trie2.RootHash())
	if strings.Compare(strRoot, hashHex) != 0 {
		t.Errorf("exp %s got %s", hashHex, strRoot)
	}
}

func updateString(trie trie.Mutable, k, v string) {
	trie.Set([]byte(k), []byte(v))
}

func deleteString(trie trie.Mutable, k string) {
	trie.Delete([]byte(k))
}
