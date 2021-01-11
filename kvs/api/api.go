// Key value store api.

package api

import (
	"errors"
	"sync"
)

var store = struct {
	sync.RWMutex
	m map[string]string
}{m: make(map[string]string)}

// Put the value into the key.
func Put(key, value string) error {
	store.Lock()
	store.m[key] = value
	store.Unlock()
	return nil
}

//ErrorNoSuchKey error value indicating key does not exist.
var ErrorNoSuchKey = errors.New("no such key")

// Get the value for a key. Returns empty string and error in case
// key does not exist.
func Get(key string) (string, error) {
	store.RLock()
	value, ok := store.m[key]
	store.RUnlock()

	if !ok {
		return "", ErrorNoSuchKey
	}

	return value, nil
}

// Delete the key.
func Delete(key string) error {
	store.Lock()
	delete(store.m, key)
	store.Unlock()
	return nil
}
