package service

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceContainerGetAndSet(t *testing.T) {
	type T1 struct{ val int }
	type T2 struct{ val float64 }

	container := NewServiceContainer()
	container.Set("a", &T1{10})
	container.Set("b", &T2{3.14})
	container.Set("c", &T1{25})

	value1, err1 := container.Get("a")
	require.Nil(t, err1)
	assert.Equal(t, &T1{10}, value1)

	value2, err2 := container.Get("b")
	require.Nil(t, err2)
	assert.Equal(t, &T2{3.14}, value2)

	value3, err3 := container.Get("c")
	require.Nil(t, err3)
	assert.Equal(t, &T1{25}, value3)
}

func TestServiceContainerInject(t *testing.T) {
	type T1 struct{ val int }
	type T2 struct {
		Value *T1 `service:"value"`
	}

	container := NewServiceContainer()
	container.Set("value", &T1{42})
	obj := &T2{}
	err := container.Inject(obj)
	require.Nil(t, err)
	assert.Equal(t, 42, obj.Value.val)
}

func TestServiceContainerInjectNonPointer(t *testing.T) {
	type T1 struct{ val int }
	type T2 struct {
		Value T1 `service:"value"`
	}

	container := NewServiceContainer()
	container.Set("value", T1{42})
	obj := &T2{}
	err := container.Inject(obj)
	require.Nil(t, err)
	assert.Equal(t, 42, obj.Value.val)
}

func TestServiceContainerInjectAnonymous(t *testing.T) {
	type T1 struct{ val int }
	type T2 struct {
		Value *T1 `service:"value"`
	}
	type T3 struct{ *T2 }

	container := NewServiceContainer()
	container.Set("value", &T1{42})
	obj := &T3{&T2{}}
	err := container.Inject(obj)
	require.Nil(t, err)
	assert.Equal(t, 42, obj.Value.val)
}

func TestServiceContainerInjectAnonymousZeroValue(t *testing.T) {
	type T1 struct{ val int }
	type T2 struct {
		Value *T1 `service:"value"`
	}
	type T3 struct{ *T2 }

	container := NewServiceContainer()
	container.Set("value", &T1{42})
	obj := &T3{} // not &T3{&T2{}}
	err := container.Inject(obj)
	require.Nil(t, err)
	assert.Equal(t, 42, obj.Value.val)
}

func TestServiceContainerInjectAnonymousNonPointer(t *testing.T) {
	type T1 struct{ val int }
	type T2 struct {
		Value *T1 `service:"value"`
	}
	type T3 struct{ T2 }

	container := NewServiceContainer()
	container.Set("value", &T1{42})
	obj := &T3{}
	err := container.Inject(obj)
	require.Nil(t, err)
	assert.Equal(t, 42, obj.Value.val)
}

func TestServiceContainerInjectAnonymousZeroValueNoServiceTags(t *testing.T) {
	type T1 struct{ val int }
	type T2 struct{ *T1 }

	container := NewServiceContainer()
	container.Set("value", &T1{42})
	obj := &T2{}
	err := container.Inject(obj)
	require.Nil(t, err)
	assert.Nil(t, obj.T1)
}

func TestServiceContainerInjectAnonymousUnexported(t *testing.T) {
	type T1 struct{ val int }
	type t2 struct {
		Value *T1 `service:"value"`
	}
	type T3 struct{ *t2 }

	container := NewServiceContainer()
	container.Set("value", &T1{42})
	obj := &T3{&t2{}}
	err := container.Inject(obj)
	require.Nil(t, err)
	assert.Nil(t, obj.t2.Value)
}

func TestServiceContainerInjectNonStruct(t *testing.T) {
	container := NewServiceContainer()
	obj := func() error { return nil }
	err := container.Inject(obj)
	require.Nil(t, err)
}

func TestServiceContainerInjectMissingService(t *testing.T) {
	type T1 struct{ val int }
	type T2 struct {
		Value *T1 `service:"value"`
	}

	container := NewServiceContainer()
	obj := &T2{}
	err := container.Inject(obj)
	assert.EqualError(t, err, `no service registered to key "value"`)
}

func TestServiceContainerInjectBadType(t *testing.T) {
	type T1 struct{ val int }
	type T2 struct {
		Value *T1 `service:"value"`
	}
	type T3 struct{ val float64 }

	container := NewServiceContainer()
	container.Set("value", &T3{3.14})
	obj := &T2{}
	err := container.Inject(obj)
	assert.EqualError(t, err, "field 'Value' cannot be assigned a value of type *service.T3")
}

func TestServiceContainerInjectNil(t *testing.T) {
	type T1 struct{ val int }
	type T2 struct {
		Value T1 `service:"value"`
	}

	container := NewServiceContainer()
	container.Set("value", nil)
	obj := &T2{}
	err := container.Inject(obj)
	assert.EqualError(t, err, "field 'Value' cannot be assigned a value of type nil")
}

func TestServiceContainerInjectOptional(t *testing.T) {
	type T1 struct{ val int }
	type T2 struct {
		Value *T1 `service:"value" optional:"true"`
	}

	container := NewServiceContainer()
	obj := &T2{}
	err := container.Inject(obj)
	require.Nil(t, err)
	require.Nil(t, obj.Value)

	container.Set("value", &T1{42})
	err = container.Inject(obj)
	require.Nil(t, err)
	assert.Equal(t, 42, obj.Value.val)
}

func TestServiceContainerInjectBadOptional(t *testing.T) {
	type T1 struct{ val int }
	type T2 struct {
		Value *T1 `service:"value" optional:"yup"`
	}

	container := NewServiceContainer()
	obj := &T2{}
	err := container.Inject(obj)
	assert.EqualError(t, err, "field 'Value' has an invalid optional tag")
}

func TestServiceContainerUnsettableFields(t *testing.T) {
	type T1 struct{ val int }
	type T2 struct {
		value *T1 `service:"value"`
	}

	container := NewServiceContainer()
	container.Set("value", &T1{42})
	err := container.Inject(&T2{})
	assert.EqualError(t, err, "field 'value' can not be set - it may be unexported")
}

func TestServiceContainerPostInject(t *testing.T) {
	container := NewServiceContainer()
	obj := &testPostInjectProcess{}
	container.Set("value", &TI{42})
	err := container.Inject(obj)
	require.Nil(t, err)
	assert.Equal(t, 42.0, obj.FValue.val)
}

type TI struct{ val int }
type TF struct{ val float64 }

type testPostInjectProcess struct {
	IValue *TI `service:"value"`
	FValue *TF
}

func (p *testPostInjectProcess) PostInject() error {
	p.FValue = &TF{float64(p.IValue.val)}
	return nil
}

func TestServiceContainerPostInjectChain(t *testing.T) {
	container := NewServiceContainer()
	obj := &testPostInjectProcessParent{}
	process := &testPostInjectProcess{}

	container.Set("value", &TI{42})
	container.Set("process", process)
	container.Set("services", container)

	err := container.Inject(obj)
	require.Nil(t, err)
	assert.Equal(t, 42.0, process.FValue.val)
}

type testPostInjectProcessParent struct {
	Services ServiceContainer       `service:"services"`
	Child    *testPostInjectProcess `service:"process"`
}

func (p *testPostInjectProcessParent) PostInject() error {
	return p.Services.Inject(p.Child)
}

func TestServiceContainerPostInjectError(t *testing.T) {
	container := NewServiceContainer()
	obj := &testPostInjectProcessError{}
	err := container.Inject(obj)
	assert.EqualError(t, err, "oops")
}

type testPostInjectProcessError struct{}

func (p *testPostInjectProcessError) PostInject() error {
	return fmt.Errorf("oops")
}

func TestServiceContainerDuplicateRegistration(t *testing.T) {
	container := NewServiceContainer()
	err1 := container.Set("dup", struct{}{})
	err2 := container.Set("dup", struct{}{})
	require.Nil(t, err1)
	assert.EqualError(t, err2, `duplicate service key "dup"`)
}

func TestServiceContainerGetUnregisteredKey(t *testing.T) {
	container := NewServiceContainer()
	_, err := container.Get("unregistered")
	assert.EqualError(t, err, `no service registered to key "unregistered"`)
}
