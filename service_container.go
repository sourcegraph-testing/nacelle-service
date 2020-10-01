package service

import (
	"fmt"
	"sync"
)

// ServiceContainer is a wrapper around services indexed by a unique
// name. Services can be retrieved by name, or injected into a struct
// by reading tagged fields.
type ServiceContainer interface {
	// Get retrieves the service registered to the given key. It is an
	// error for a service not to be registered to this key.
	Get(key string) (interface{}, error)

	// MustGet calls Get and panics on error.
	MustGet(service string) interface{}

	// Set registers a service with the given key. It is an error for
	// a service to already be registered to this key.
	Set(key string, service interface{}) error

	// MustSet calls Set and panics on error.
	MustSet(service string, value interface{})

	// Inject will attempt to populate the given type with values from
	// the service container based on the value's struct tags. An error
	// may occur if a service has not been registered, a service has a
	// different type than expected, or struct tags are malformed.
	Inject(obj interface{}) error
}

type serviceContainer struct {
	services map[string]interface{}
	mutex    sync.RWMutex
}

// NewServiceContainer creates an empty service container.
func NewServiceContainer() ServiceContainer {
	return &serviceContainer{
		services: map[string]interface{}{},
	}
}

// Get retrieves the service registered to the given key. It is an
// error for a service not to be registered to this key.
func (c *serviceContainer) Get(key string) (interface{}, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	service, ok := c.services[key]
	if !ok {
		return nil, fmt.Errorf("no service registered to key `%s`", key)
	}

	return service, nil
}

// MustGet calls Get and panics on error.
func (c *serviceContainer) MustGet(service string) interface{} {
	value, err := c.Get(service)
	if err != nil {
		panic(err.Error())
	}

	return value
}

// Set registers a service with the given key. It is an error for
// a service to already be registered to this key.
func (c *serviceContainer) Set(key string, service interface{}) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, ok := c.services[key]; ok {
		return fmt.Errorf("duplicate service key `%s`", key)
	}

	c.services[key] = service
	return nil
}

// MustSet calls Set and panics on error.
func (c *serviceContainer) MustSet(service string, value interface{}) {
	if err := c.Set(service, value); err != nil {
		panic(err.Error())
	}
}

// Inject will attempt to populate the given type with values from
// the service container based on the value's struct tags. An error
// may occur if a service has not been registered, a service has a
// different type than expected, or struct tags are malformed.
func (c *serviceContainer) Inject(obj interface{}) error {
	_, err := inject(c, obj, nil, nil)
	return err
}
