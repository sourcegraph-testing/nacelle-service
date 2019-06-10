# Bussard

[![GoDoc](https://godoc.org/github.com/go-nacelle/service?status.svg)](https://godoc.org/github.com/go-nacelle/service)
[![Build Status](https://secure.travis-ci.org/go-nacelle/service.png)](http://travis-ci.org/go-nacelle/service)
[![Maintainability](https://api.codeclimate.com/v1/badges/5f7ceba80716e77fe9fe/maintainability)](https://codeclimate.com/github/go-nacelle/service/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/5f7ceba80716e77fe9fe/test_coverage)](https://codeclimate.com/github/go-nacelle/service/test_coverage)

[![CircleCI](https://circleci.com/gh/go-nacelle/service.svg?style=svg)](https://circleci.com/gh/go-nacelle/service)
[![Coverage Status](https://coveralls.io/repos/github/go-nacelle/service/badge.svg?branch=master)](https://coveralls.io/github/go-nacelle/service?branch=master)

Bussard is a service container and dependency injection tool for Golang.

## Overview

Bussard provides a `ServiceContainer` to which a service (any value) can
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

## License

Copyright (c) 2018 Eric Fritz

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
