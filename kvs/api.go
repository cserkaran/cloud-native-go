// Key value store api.

package api

import "errors"

var store = make(map[string]string)

// Put the value into the key.
func Put(key, value string) error {
	store[key] = value
	return nil
}

//ErrorNoSuchKey error value indicating key does not exist.
var ErrorNoSuchKey = errors.New("no such key")

// Get the value for a key. Returns empty string and error in case
// key does not exist.
func Get(key string) (string, error) {
	value, ok := store[key]

	if !ok {
		return "", ErrorNoSuchKey
	}

	return value, nil
}

// Delete the key.
func Delete(key string) error {
	delete(store, key)
	return nil
}
