# Nacelle Service Container [![GoDoc](https://godoc.org/github.com/go-nacelle/service?status.svg)](https://godoc.org/github.com/go-nacelle/service) [![CircleCI](https://circleci.com/gh/go-nacelle/service.svg?style=svg)](https://circleci.com/gh/go-nacelle/service) [![Coverage Status](https://coveralls.io/repos/github/go-nacelle/service/badge.svg?branch=master)](https://coveralls.io/github/go-nacelle/service?branch=master)

Service container and dependency injection for [nacelle](https://github.com/go-nacelle/nacelle).

---

A **service** is a value that can injected into one or more structs. Generally, a service is going to be something that controls access to shared state (e.g. a database connection or in-memory cache). A **service container** is a collection of named services. Services are registered into a container and then the container is used to injected services into consumer objects.

### Usage

Basic usage is a generic key-value store. Once a service container is created, its `Get` and `Set` method can be used to retrieve and register services by name, respectively.

```go
services := NewServiceContainer()
if err := services.Set("example", &SomeService{}); err != nil {
    // handle error
}

service, err : services.Get("example").(*SomeService)
if err != nil {
    // handle error
}
```

The `Get` method fails when no such service is registered, and the `Set` method fails when a service name is reused. The `Get` and `Set` methods have analogous `MustGet` and `MustSet` methods which panic on error.

The `Get` method returns a bare interface object, as the service container holds services of heterogeneous type. The standard usage of a service container is to **inject** service instances into tagged structs, which does type conversions for you.

```go
type Consumer struct {
    Service *SomeExample `service:"example"`
}

consumer := &Consumer{}
if err := container.Inject(consumer); err != nil {
    // handle error
}
```

The `Inject` method fails when a consumer asks for an unregistered service or for a service with the wrong target type. Services can be tagged as optional (e.g. `service:"example" optional:"true"`) which will silence the later kind of error. Tagged fields must be exported.

### Post Injection Hook

After successful injection to a struct, the method named `PostInject` will be called if it is defined. This allows initialization behavior that relies on the presence of injected services to be run as soon as possible.

```go
func (c *Consumer) PostInject() error {
    return c.Service.PrepFor("consumer")
}
```

### Anonymous Structs

Injection also works on structs containing composite fields. The following example successfully assigns the registered value to the field `Child.Base.Service`.

```go
type Base struct {
    Service *SomeExample `service:"example"`
}

type Child struct {
    *Base
}

child := &Child{}
if err := container.Inject(child); err != nil {
    // handle error
}
```

### Recursive Injection

It should be noted that injection does not work **recursively** -- the procedure does not look into the values of non-anonymous fields. If this behavior is needed, it can be performed during a post-injection hook. The following example demonstrates this behavior and assumes that the service container is registered to itself under the name `services`.

```go
type RecursiveInjectionConsumer struct {
    Services service.ServiceContainer `service:"services"`
    FieldA   *A
    FieldB   *B
    FieldC   *C
}

func (c *RecursiveInjectionConsumer) PostInject() error {
    for _, field := range []interface{}{c.FieldA, c.FieldB, c.FieldC} {
        if err := c.Services.Inject(field); err != nil {
            return err
        }
    }

    return nil
}
```
