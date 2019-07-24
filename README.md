# Nacelle Service Container [![GoDoc](https://godoc.org/github.com/go-nacelle/service?status.svg)](https://godoc.org/github.com/go-nacelle/service) [![CircleCI](https://circleci.com/gh/go-nacelle/service.svg?style=svg)](https://circleci.com/gh/go-nacelle/service) [![Coverage Status](https://coveralls.io/repos/github/go-nacelle/service/badge.svg?branch=master)](https://coveralls.io/github/go-nacelle/service?branch=master)

Service container and dependency injection for [nacelle](https://nacelle.dev).

---

A **service container** is a collection of objects which are constructed separately from their consumers. This pattern allows for a greater separation of concerns, where consumers care only about a particular concrete or interface type, but do not care about their configuration, construction, or initialization. This separation also allows multiple consumers for the same service which does not need to be initialized multiple times (e.g. a database connection or an in-memory cache layer). Service injection is performed on [initializers and processes](https://nacelle.dev/docs/core/process) during application startup automatically.

You can see an additional example of service injection in the [example repository](https://github.com/go-nacelle/example), specifically the [worker spec](https://github.com/go-nacelle/example/blob/843979aaa86786784a1ca3646e8d0d1f69e29c65/cmd/worker/worker_spec.go#L15) definition. In this project, the `Conn` and `PubSubConn` services are created by application-defined initializers [here](https://github.com/go-nacelle/example/blob/843979aaa86786784a1ca3646e8d0d1f69e29c65/internal/redis_initializer.go#L28) and [here](https://github.com/go-nacelle/example/blob/843979aaa86786784a1ca3646e8d0d1f69e29c65/internal/pubsub_initializer.go#L32).

### Registration

A concrete service can be registered to the service container with a unique name by which it can later be retrieved. The `Set` method fails when a service name is reused. There is also an analogous `MustSet` method that panics on error.

```go
func Init(services nacelle.ServiceContainer) error {
    example := &Example{}
    if err := services.Set("example", example); err != nil {
        return err
    }

    // ...
}
```

The [logger](https://nacelle.dev/docs/core/log) (under the name `logger`), the [health tracker](https://nacelle.dev/docs/core/process#tracking-process-health) (under the name `health`), and the service container itself (under the name `services`) are available in all applications using the nacelle [bootstrapper](https://nacelle.dev/docs/core).

### Retrieval

A service can be retrieved from the service container by the name with which it is registered. However, this returns a bare interface object and requires the consumer of the service to do a type-check and cast.

Instead, the recommended way to consume dependent services is to **inject** them into a struct decorated with tagged fields. This does the proper type conversion for you.

```go
type Consumer struct {
    Service *SomeExample `service:"example"`
}

consumer := &Consumer{}
if err := services.Inject(consumer); err != nil {
    // handle error
}
```

The `Inject` method fails when a consumer asks for an unregistered service or for a service with the wrong target type. Services can be tagged as optional (e.g. `service:"example" optional:"true"`) which will silence the later class of errors. Tagged fields must be exported.

#### Post Injection Hook

If additional behavior is necessary after services are available to a consumer struct (e.g. running injection on the elements of a dynamic slice or map or cross-linking dependencies), the method `PostInject` can be implemented. This method, if it is defined, is invoked immediately after successful injection.

```go
func (c *Consumer) PostInject() error {
    return c.Service.PrepFor("consumer")
}
```

#### Anonymous Structs

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

#### Recursive Injection

It should be noted that injection does not work **recursively**. The procedure does not look into the values of non-anonymous fields. If this behavior is needed, it can be performed during a post-injection hook. The following example demonstrates this behavior.

```go
type RecursiveInjectionConsumer struct {
    Services service.ServiceContainer `service:"services"`
    FieldA   *A
    FieldB   *B
    FieldC   *C
}

func (c *RecursiveInjectionConsumer) PostInject() error {
    fields := []interface{}{
        c.FieldA,
        c.FieldB,
        c.FieldC,
    }

    for _, field := range fields {
        if err := c.Services.Inject(field); err != nil {
            return err
        }
    }

    return nil
}
```
