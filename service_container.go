package service

import (
	"fmt"
	"reflect"
	"strconv"
)

type (
	// ServiceContainer is a wrapper around services indexed by a unique
	// name. Services can be retrieved by name, or injected into a struct
	// by reading tagged fields.
	ServiceContainer interface {
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

	// PostInject is a marker interface for injectable objects which should
	// perform some action after injection of services.
	PostInject interface {
		PostInject() error
	}

	serviceContainer struct {
		services map[string]interface{}
	}
)

const (
	serviceTag  = "service"
	optionalTag = "optional"
)

// NewServiceContainer creates an empty service container.
func NewServiceContainer() ServiceContainer {
	return &serviceContainer{
		services: map[string]interface{}{},
	}
}

// Get retrieves the service registered to the given key. It is an
// error for a service not to be registered to this key.
func (c *serviceContainer) Get(key string) (interface{}, error) {
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
	_, err := c.inject(obj, nil, nil)
	return err
}

func (c *serviceContainer) inject(obj interface{}, root *reflect.Value, baseIndexPath []int) (bool, error) {
	var (
		ov = reflect.ValueOf(obj)
		oi = reflect.Indirect(ov)
		ot = oi.Type()
	)

	if root == nil {
		root = &oi
	}

	if oi.Kind() != reflect.Struct {
		return false, nil
	}

	hasTag := false
	for i := 0; i < ot.NumField(); i++ {
		indexPath := make([]int, len(baseIndexPath))
		copy(indexPath, baseIndexPath)
		indexPath = append(indexPath, i)

		var (
			fieldType   = ot.Field(i)
			fieldValue  = (*root).FieldByIndex(indexPath)
			serviceTag  = fieldType.Tag.Get(serviceTag)
			optionalTag = fieldType.Tag.Get(optionalTag)
		)

		if fieldType.Anonymous {
			if !fieldValue.CanSet() {
				continue
			}

			wasZeroValue := false
			if !reflect.Indirect(fieldValue).IsValid() {
				initializedValue := reflect.New(fieldType.Type.Elem())
				fieldValue.Set(initializedValue)
				fieldValue = initializedValue
				wasZeroValue = true
			}

			anonymousFieldHasTag, err := c.inject(fieldValue.Interface(), root, indexPath)
			if err != nil {
				return false, err
			}

			if anonymousFieldHasTag {
				hasTag = true
			} else if wasZeroValue {
				zeroValue := reflect.Zero(fieldType.Type)
				fieldValue = (*root).FieldByIndex(indexPath)
				fieldValue.Set(zeroValue)
			}

			continue
		}

		if serviceTag == "" {
			continue
		}

		hasTag = true

		if err := loadServiceField(c, fieldType, fieldValue, serviceTag, optionalTag); err != nil {
			return false, err
		}
	}

	if pi, ok := obj.(PostInject); ok {
		if err := pi.PostInject(); err != nil {
			return false, err
		}
	}

	return hasTag, nil
}

func loadServiceField(container *serviceContainer, fieldType reflect.StructField, fieldValue reflect.Value, serviceTag, optionalTag string) error {
	if !fieldValue.IsValid() {
		return fmt.Errorf("field '%s' is invalid", fieldType.Name)
	}

	if !fieldValue.CanSet() {
		return fmt.Errorf("field '%s' can not be set - it may be unexported", fieldType.Name)
	}

	value, err := container.Get(serviceTag)
	if err != nil {
		if optionalTag != "" {
			val, err := strconv.ParseBool(optionalTag)
			if err != nil {
				return fmt.Errorf("field '%s' has an invalid optional tag", fieldType.Name)
			}

			if val {
				return nil
			}
		}

		return err
	}

	var (
		targetType  = fieldValue.Type()
		targetValue = reflect.ValueOf(value)
	)

	if !targetValue.IsValid() || !targetValue.Type().ConvertibleTo(targetType) {
		return fmt.Errorf(
			"field '%s' cannot be assigned a value of type %s",
			fieldType.Name,
			getTypeName(value),
		)
	}

	fieldValue.Set(targetValue.Convert(targetType))
	return nil
}

func getTypeName(v interface{}) string {
	if v == nil {
		return "nil"
	}

	return reflect.TypeOf(v).String()
}
