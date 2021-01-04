package service

import (
	"fmt"
	"reflect"
)

// InjectableServiceKey is an optional interface for service keys. Non-string service key values
// should implement this interface if they intend to be injected via struct tags.
type InjectableServiceKey interface {
	// Tag returns the string version of the key. Two distinct key values that return the same
	// value from this method should be considered equivalent within a single service container.
	Tag() string
}

// prettyKey returns a human-readable string describing the given
// service key.
func prettyKey(key interface{}) string {
	if tag, ok := tagForKey(key); ok {
		if _, ok := key.(string); ok {
			return fmt.Sprintf(`"%s"`, key)
		}

		return fmt.Sprintf(`%s ("%s")`, reflect.TypeOf(key).Name(), tag)
	}

	return reflect.TypeOf(key).Name()
}

// tagForKey returns the string version of the given struct key value and a boolean flag indicating
// such a string's existence.
func tagForKey(key interface{}) (string, bool) {
	if k, ok := key.(string); ok {
		return k, true
	}

	if k, ok := key.(InjectableServiceKey); ok {
		return k.Tag(), true
	}

	return "", false
}
