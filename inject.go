package service

import (
	"fmt"
	"reflect"
	"strconv"
)

// Inject will attempt to populate the given type with values from the service container based on
// the value's struct tags. An error may occur if a service has not been registered, a service has
// a different type than expected, or struct tags are malformed.
func Inject(c *Container, obj interface{}) error {
	_, err := inject(c, obj, nil, nil)
	return err
}

// inject populates fields of the given struct. The root parameter should always point to the top
// of the struct object. Passing nil will set the root to be the reflected value of the given object.
// The given integer path should be the field index path to the object from the root of the struct.
// This function returns true if the struct value was updated. If the object conforms to the PostInject
// interface, its hook is called after successful injection.
func inject(c *Container, obj interface{}, root *reflect.Value, path []int) (bool, error) {
	oi := reflect.Indirect(reflect.ValueOf(obj))
	if oi.Kind() != reflect.Struct {
		return false, nil
	}

	ot := oi.Type()

	if root == nil {
		root = &oi
	}

	updated := false
	for i := 0; i < ot.NumField(); i++ {
		fieldPath := make([]int, len(path), len(path)+1)
		copy(path, path)
		fieldPath = append(fieldPath, i)

		fieldUpdated, err := injectField(c, ot.Field(i), root, fieldPath)
		if err != nil {
			return false, err
		}

		updated = updated || fieldUpdated
	}

	if pi, ok := obj.(PostInject); ok {
		if err := pi.PostInject(); err != nil {
			return false, err
		}
	}

	return updated, nil
}

const (
	serviceTag  = "service"
	optionalTag = "optional"
)

// injectField recursively sets the value of the given struct field. This uses the service struct tag
// as the service key to match in the given container. If the field is a nested anonymous struct, its
// fields are injected recursively. This function returns true if the field was updated.
func injectField(c *Container, fieldType reflect.StructField, root *reflect.Value, indexPath []int) (bool, error) {
	if fieldType.Anonymous {
		return injectAnonymousField(c, fieldType, root, indexPath)
	}

	fieldValue := (*root).FieldByIndex(indexPath)
	serviceTag := fieldType.Tag.Get(serviceTag)
	optionalTag := fieldType.Tag.Get(optionalTag)

	if serviceTag == "" {
		return false, nil
	}

	optional := false
	if optionalTag != "" {
		val, err := strconv.ParseBool(optionalTag)
		if err != nil {
			return false, fmt.Errorf("field '%s' has an invalid optional tag", fieldType.Name)
		}

		optional = val
	}

	return loadServiceField(c, fieldType, fieldValue, serviceTag, optional)
}

// injectAnonymousField sets the value of the given struct field to the recursively injected value
// for this field. If the field is unset, a zero value of the field's type will be used as a base.
// This function returns true if the struct field was updated.
func injectAnonymousField(c *Container, fieldType reflect.StructField, root *reflect.Value, indexPath []int) (bool, error) {
	fieldValue := (*root).FieldByIndex(indexPath)
	if !fieldValue.CanSet() {
		return false, nil
	}

	wasZeroValue := false
	if !reflect.Indirect(fieldValue).IsValid() {
		wasZeroValue = true
		initializedValue := reflect.New(fieldType.Type.Elem())
		fieldValue.Set(initializedValue)
		fieldValue = initializedValue
	}

	anonymousFieldHasTag, err := inject(c, fieldValue.Interface(), root, indexPath)
	if err != nil {
		return false, err
	}

	if !anonymousFieldHasTag && wasZeroValue {
		zeroValue := reflect.Zero(fieldType.Type)
		fieldValue = (*root).FieldByIndex(indexPath)
		fieldValue.Set(zeroValue)
	}

	return true, nil
}

// loadServiceField sets the value of the given struct field to the value of the service registered to
// the given service key in the given container. This function returns true if the field was updated.
func loadServiceField(c *Container, fieldType reflect.StructField, fieldValue reflect.Value, serviceTag string, optional bool) (bool, error) {
	if !fieldValue.IsValid() {
		return false, fmt.Errorf("field '%s' is invalid", fieldType.Name)
	}

	if !fieldValue.CanSet() {
		return false, fmt.Errorf("field '%s' can not be set - it may be unexported", fieldType.Name)
	}

	value, err := c.Get(serviceTag)
	if err != nil {
		if optional {
			return false, nil
		}

		return false, err
	}

	targetType := fieldValue.Type()
	targetValue := reflect.ValueOf(value)

	if !targetValue.IsValid() || !targetValue.Type().ConvertibleTo(targetType) {
		var typeName string
		if value == nil {
			typeName = "nil"
		} else {
			typeName = reflect.TypeOf(value).String()
		}

		return false, fmt.Errorf("field '%s' cannot be assigned a value of type %s", fieldType.Name, typeName)
	}

	fieldValue.Set(targetValue.Convert(targetType))
	return true, nil
}
