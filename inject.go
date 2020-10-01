package service

import (
	"fmt"
	"reflect"
	"strconv"
)

// ServiceGetter is a subset of a ServiceContainer that only supports the
// retrieval of a registered service by name.
type ServiceGetter interface {
	// Get retrieves the service registered to the given key. It is an
	// error for a service not to be registered to this key.
	Get(key string) (interface{}, error)
}

// PostInject is a marker interface for injectable objects which should
// perform some action after injection of services.
type PostInject interface {
	PostInject() error
}

const (
	serviceTag  = "service"
	optionalTag = "optional"
)

func inject(c ServiceGetter, obj interface{}, root *reflect.Value, baseIndexPath []int) (bool, error) {
	ov := reflect.ValueOf(obj)
	oi := reflect.Indirect(ov)
	ot := oi.Type()

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

		fieldType := ot.Field(i)
		fieldValue := (*root).FieldByIndex(indexPath)
		serviceTag := fieldType.Tag.Get(serviceTag)
		optionalTag := fieldType.Tag.Get(optionalTag)

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

			anonymousFieldHasTag, err := inject(c, fieldValue.Interface(), root, indexPath)
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

func loadServiceField(container ServiceGetter, fieldType reflect.StructField, fieldValue reflect.Value, serviceTag, optionalTag string) error {
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

	targetType := fieldValue.Type()
	targetValue := reflect.ValueOf(value)

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
