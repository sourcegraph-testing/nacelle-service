package service

type overlayContainer struct {
	base     ServiceContainer
	services map[interface{}]interface{}
}

// Overlay wraps the given service container with an immutable map of
// services. Calling Get on the resulting service container will return
// a service from the overlay map, then will fall back to the wrapped
// service container. Similarly, Inject will favor services from the
// overlay map.
//
// This allows a user to re-assign services in the container for a specific
// specialized code path. This can be used, for example, to inject a logger
// with context for the current request or task to a short-lived handler.
//
// Calling Set will modify the wrapped container directly.
func Overlay(base ServiceContainer, services map[interface{}]interface{}) ServiceContainer {
	return &overlayContainer{
		base:     base,
		services: services,
	}
}

// Get retrieves the service registered to the given key. It is an
// error for a service not to be registered to this key.
func (c *overlayContainer) Get(key interface{}) (interface{}, error) {
	if service, ok := c.services[key]; ok {
		return service, nil
	}

	return c.base.Get(key)
}

// Set registers a service with the given key. It is an error for
// a service to already be registered to this key.
func (c *overlayContainer) Set(key interface{}, service interface{}) error {
	return c.base.Set(key, service)
}

// Inject will attempt to populate the given type with values from
// the service container based on the value's struct tags. An error
// may occur if a service has not been registered, a service has a
// different type than expected, or struct tags are malformed.
func (c *overlayContainer) Inject(obj interface{}) error {
	_, err := inject(c, obj, nil, nil)
	return err
}
