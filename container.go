package service

import (
	"fmt"
	"sync"
)

// Container is a collection of services retrievable by a unique service key value.
type Container struct {
	services  map[interface{}]interface{}
	keysByTag map[string]interface{}
	parent    *Container
	mutex     sync.RWMutex
}

// New creates an empty service container.
func New() *Container {
	return &Container{
		services:  map[interface{}]interface{}{},
		keysByTag: map[string]interface{}{},
	}
}

// Get retrieves the service registered to the given key. It is an error for a service not
// to be registered to this key.
func (c *Container) Get(key interface{}) (interface{}, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	// Service exists under key
	if service, ok := c.services[key]; ok {
		return service, nil
	}
	if tag, ok := tagForKey(key); ok {
		if key, ok := c.keysByTag[tag]; ok {
			// Service exists under key with same tag
			if service, ok := c.services[key]; ok {
				return service, nil
			}
		}
	}

	if c.parent != nil {
		// Check parent layers
		return c.parent.Get(key)
	}

	return nil, fmt.Errorf("no service registered to key %s", prettyKey(key))
}

// Set registers a service with the given key. It is an error for a service to already be
// registered to this key (or a key with the same tag, see InjectableServiceKey).
func (c *Container) Set(key, service interface{}) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Service exists under key
	if _, ok := c.services[key]; ok {
		return fmt.Errorf(`duplicate service key %s`, prettyKey(key))
	}

	tag, ok := tagForKey(key)
	if ok {
		// Service exists under key with same tag
		if _, ok := c.keysByTag[tag]; ok {
			return fmt.Errorf(`duplicate service key %s`, prettyKey(key))
		}
	}

	if c.parent != nil {
		// Delegate to parent if we're not the root
		return c.parent.Set(key, service)
	}

	// We're the root, update both maps
	c.services[key] = service
	if ok {
		c.keysByTag[tag] = key
	}

	return nil
}

// WithValues returns a copy of the container with the given service map overlaid on top.
// Calling  Set on the resulting container will modify the original container and any other
// containers created from this method. It is an error for the given map to contain two keys
// that resolve to the same tag (see InjectableServiceKey).
func (c *Container) WithValues(services map[interface{}]interface{}) (*Container, error) {
	c2 := New()
	for k, v := range services {
		if err := c2.Set(k, v); err != nil {
			return nil, err
		}
	}

	c2.parent = c
	return c2, nil
}
