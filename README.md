# Nacelle Service Container [![GoDoc](https://godoc.org/github.com/go-nacelle/service?status.svg)](https://godoc.org/github.com/go-nacelle/service) [![CircleCI](https://circleci.com/gh/go-nacelle/service.svg?style=svg)](https://circleci.com/gh/go-nacelle/service) [![Coverage Status](https://coveralls.io/repos/github/go-nacelle/service/badge.svg?branch=master)](https://coveralls.io/github/go-nacelle/service?branch=master)

Service container and dependency injection for [nacelle](https://github.com/go-nacelle/nacelle).

---

This library provides a `ServiceContainer` to which a service (any value) can
be registered with a unique name. Users of that service can pull a reference
to the service from the container by using the same name.

Usage of this library can improve code which must thread references to objects
which are required either frequently, or somewhere deep within the application.
Instead, a single populated container instance can be passed and clients can
simply pull their own dependencies from it without requiring the calling code
to worry about initialization dependencies.

## Usage

Basic usage is very simple. Simply create a container, then get and set objects
into it. Names must be unique, so setting an object requires the name is not
already taken, and getting an object requires that the object has already been
put into the container.

```go
container := NewServiceContainer()
if err := container.Set("name", &NewService()); err != nil {
    // ...
}

// Later that day...
service, err := container.Get("name").(*Service)
if err != nil {
    // ...
}
```

The `Get` and `Set` methods have analogous `MustGet` and `MustSet` methods
which panic on error.

Standard usage *injects* service instances into tagged structs. This should
be the default way in which a client of a service receives an instance.

The following struct tags four fields with service dependencies. When the
container injects services into this instance, it will look up services with
the given name and attempt to cast them to the appropriate type. It is an error
to request a service that has not been registered (unless the optional tag is
present, in which case the field value remains nil), or to request a service
of a type which cannot be assigned to the field. Notice that fields with tags
must be exported.

```go
type Server struct {
    Logger         logging.Logger     `service:"logger"`
    MetricReporter metrics.Reporter   `service:"metric-reporter"`
    Cache          cache.Cache        `service:"cache"`
    SecondaryAuth  auth.Authenticator `service:"authenticator" optional:"true"`
}
```

To inject services, the code which creates the server (depending on application
flow) must simply pass the instance to inject to the container.

```go
server := &Server{}

// ...

if err := container.Inject(server); err != nil {
    // ...
}
```
