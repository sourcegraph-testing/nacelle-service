package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInject(t *testing.T) {
	type T1 struct{ val int }
	type T2 struct {
		Value *T1 `service:"value"`
	}

	container := New()
	container.Set("value", &T1{42})
	obj := &T2{}
	err := Inject(context.Background(), container, obj)
	require.Nil(t, err)
	assert.Equal(t, 42, obj.Value.val)
}

func TestInjectNonPointer(t *testing.T) {
	type T1 struct{ val int }
	type T2 struct {
		Value T1 `service:"value"`
	}

	container := New()
	container.Set("value", T1{42})
	obj := &T2{}
	err := Inject(context.Background(), container, obj)
	require.Nil(t, err)
	assert.Equal(t, 42, obj.Value.val)
}

func TestInjectAnonymous(t *testing.T) {
	type T1 struct{ val int }
	type T2 struct {
		Value *T1 `service:"value"`
	}
	type T3 struct{ *T2 }

	container := New()
	container.Set("value", &T1{42})
	obj := &T3{&T2{}}
	err := Inject(context.Background(), container, obj)
	require.Nil(t, err)
	assert.Equal(t, 42, obj.Value.val)
}

func TestInjectAnonymousZeroValue(t *testing.T) {
	type T1 struct{ val int }
	type T2 struct {
		Value *T1 `service:"value"`
	}
	type T3 struct{ *T2 }

	container := New()
	container.Set("value", &T1{42})
	obj := &T3{} // not &T3{&T2{}}
	err := Inject(context.Background(), container, obj)
	require.Nil(t, err)
	assert.Equal(t, 42, obj.Value.val)
}

func TestInjectAnonymousNonPointer(t *testing.T) {
	type T1 struct{ val int }
	type T2 struct {
		Value *T1 `service:"value"`
	}
	type T3 struct{ T2 }

	container := New()
	container.Set("value", &T1{42})
	obj := &T3{}
	err := Inject(context.Background(), container, obj)
	require.Nil(t, err)
	assert.Equal(t, 42, obj.Value.val)
}

func TestInjectAnonymousZeroValueNoServiceTags(t *testing.T) {
	type T1 struct{ val int }
	type T2 struct{ *T1 }

	container := New()
	container.Set("value", &T1{42})
	obj := &T2{}
	err := Inject(context.Background(), container, obj)
	require.Nil(t, err)
	assert.Nil(t, obj.T1)
}

func TestInjectAnonymousUnexported(t *testing.T) {
	type T1 struct{ val int }
	type t2 struct {
		Value *T1 `service:"value"`
	}
	type T3 struct{ *t2 }

	container := New()
	container.Set("value", &T1{42})
	obj := &T3{&t2{}}
	err := Inject(context.Background(), container, obj)
	require.Nil(t, err)
	assert.Nil(t, obj.t2.Value)
}

func TestInjectNonStruct(t *testing.T) {
	container := New()
	obj := func() error { return nil }
	err := Inject(context.Background(), container, obj)
	require.Nil(t, err)
}

func TestInjectMissingService(t *testing.T) {
	type T1 struct{ val int }
	type T2 struct {
		Value *T1 `service:"value"`
	}

	container := New()
	obj := &T2{}
	err := Inject(context.Background(), container, obj)
	assert.EqualError(t, err, `no service registered to key "value"`)
}

func TestInjectBadType(t *testing.T) {
	type T1 struct{ val int }
	type T2 struct {
		Value *T1 `service:"value"`
	}
	type T3 struct{ val float64 }

	container := New()
	container.Set("value", &T3{3.14})
	obj := &T2{}
	err := Inject(context.Background(), container, obj)
	assert.EqualError(t, err, "field 'Value' cannot be assigned a value of type *service.T3")
}

func TestInjectNil(t *testing.T) {
	type T1 struct{ val int }
	type T2 struct {
		Value T1 `service:"value"`
	}

	container := New()
	container.Set("value", nil)
	obj := &T2{}
	err := Inject(context.Background(), container, obj)
	assert.EqualError(t, err, "field 'Value' cannot be assigned a value of type nil")
}

func TestInjectOptional(t *testing.T) {
	type T1 struct{ val int }
	type T2 struct {
		Value *T1 `service:"value" optional:"true"`
	}

	container := New()
	obj := &T2{}
	err := Inject(context.Background(), container, obj)
	require.Nil(t, err)
	require.Nil(t, obj.Value)

	container.Set("value", &T1{42})
	err = Inject(context.Background(), container, obj)
	require.Nil(t, err)
	assert.Equal(t, 42, obj.Value.val)
}

func TestInjectBadOptional(t *testing.T) {
	type T1 struct{ val int }
	type T2 struct {
		Value *T1 `service:"value" optional:"yup"`
	}

	container := New()
	obj := &T2{}
	err := Inject(context.Background(), container, obj)
	assert.EqualError(t, err, "field 'Value' has an invalid optional tag")
}

func TestContainerUnsettableFields(t *testing.T) {
	type T1 struct{ val int }
	type T2 struct {
		value *T1 `service:"value"`
	}

	container := New()
	container.Set("value", &T1{42})
	err := Inject(context.Background(), container, &T2{})
	assert.EqualError(t, err, "field 'value' can not be set - it may be unexported")
}
