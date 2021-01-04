package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContainerGetAndSet(t *testing.T) {
	type T1 struct{ val int }
	type T2 struct{ val float64 }

	container := New()
	container.Set("a", &T1{10})
	container.Set("b", &T2{3.14})
	container.Set("c", &T1{25})

	assertValue(t, container, "a", &T1{10})
	assertValue(t, container, "b", &T2{3.14})
	assertValue(t, container, "c", &T1{25})
}

func TestContainerGetAndSetNonStringKey(t *testing.T) {
	type T struct{ val int }

	key1 := testKey3{"foo"}
	key2 := testKey3{"bar"}

	container := New()
	require.Nil(t, container.Set(key1, &T{10}))
	require.Nil(t, container.Set(key2, &T{25}))

	assertValue(t, container, key1, &T{10})
	assertValue(t, container, key2, &T{25})
}

func TestContainerGetAndSetInjectableServiceKey(t *testing.T) {
	type T struct{ val int }

	key1 := testKey1{"foo"}
	key2 := testKey1{"bar"}

	container := New()
	require.Nil(t, container.Set(key1, &T{10}))
	require.Nil(t, container.Set(key2, &T{25}))

	assertValue(t, container, key1, &T{10})
	assertValue(t, container, key2, &T{25})
	assertValue(t, container, "foo", &T{10})
	assertValue(t, container, "bar", &T{25})
}

func TestContainerGetUnknownKey(t *testing.T) {
	container := New()
	_, err := container.Get("unregistered")
	assert.EqualError(t, err, `no service registered to key "unregistered"`)
}

func TestContainerSetDuplicateKey(t *testing.T) {
	container := New()
	err1 := container.Set("dup", struct{}{})
	err2 := container.Set("dup", struct{}{})
	require.Nil(t, err1)
	assert.EqualError(t, err2, `duplicate service key "dup"`)
}

func TestContainerGetAndSetDuplicateNonStringKey(t *testing.T) {
	key1 := testKey3{"dup"}
	key2 := testKey3{"dup"}

	container := New()
	require.Nil(t, container.Set(key1, struct{}{}))
	assert.EqualError(t, container.Set(key2, struct{}{}), `duplicate service key testKey3`)
}

func TestContainerGetAndSetDuplicateInjectableServiceKey(t *testing.T) {
	key1 := testKey1{"dup"}
	key2 := testKey2{"dup"}

	container := New()
	require.Nil(t, container.Set(key1, struct{}{}))
	assert.EqualError(t, container.Set(key2, struct{}{}), `duplicate service key testKey2 ("dup")`)
}

func TestContainerWithValues(t *testing.T) {
	testContainerWithValues(t, "a", "b", "c", "d")
}

func TestContainerWithValuesNonStringKey(t *testing.T) {
	testContainerWithValues(t, testKey3{"a"}, testKey3{"b"}, testKey3{"c"}, testKey3{"d"})
}

func TestContainerWithValuesInjectableServiceKey(t *testing.T) {
	testContainerWithValues(t, testKey1{"a"}, testKey1{"b"}, testKey1{"c"}, testKey1{"d"})
}

func testContainerWithValues(t *testing.T, key1, key2, key3, key4 interface{}) {
	type T struct{ val int }

	container1 := New()
	container1.Set(key1, &T{10})
	container1.Set(key2, &T{20})
	container1.Set(key3, &T{30})

	container2, err := container1.WithValues(map[interface{}]interface{}{
		key1: &T{25},
		key4: &T{50},
	})
	require.Nil(t, err)

	container3, err := container1.WithValues(map[interface{}]interface{}{
		key2: &T{50},
		key4: &T{75},
	})
	require.Nil(t, err)

	assertValue(t, container1, key1, &T{10})
	assertValue(t, container1, key2, &T{20})
	assertValue(t, container1, key3, &T{30})

	assertValue(t, container2, key1, &T{25})
	assertValue(t, container2, key2, &T{20})
	assertValue(t, container2, key3, &T{30})
	assertValue(t, container2, key4, &T{50})

	assertValue(t, container3, key1, &T{10})
	assertValue(t, container3, key2, &T{50})
	assertValue(t, container3, key3, &T{30})
	assertValue(t, container3, key4, &T{75})
}

func TestContainerWithValuesDuplicateInjectableServiceKey(t *testing.T) {
	key1 := testKey1{"dup"}
	key2 := testKey2{"dup"}

	// One will (non-deterministically) clash with the other
	_, err := New().WithValues(map[interface{}]interface{}{
		key1: struct{}{},
		key2: struct{}{},
	})

	expectedErrors := []string{
		`duplicate service key testKey1 ("dup")`,
		`duplicate service key testKey2 ("dup")`,
	}
	require.NotNil(t, err)
	assert.Contains(t, expectedErrors, err.Error())
}

func TestContainerWithValuesSet(t *testing.T) {
	type T struct{ val int }

	container1 := New()
	container1.Set("a", &T{10})
	container1.Set("b", &T{20})
	container1.Set("c", &T{30})

	container2, err := container1.WithValues(map[interface{}]interface{}{
		"a": &T{25},
		"d": &T{50},
	})
	require.Nil(t, err)

	require.Nil(t, container1.Set("d", &T{75}))
	require.Nil(t, container2.Set("e", &T{75}))

	assertValue(t, container1, "a", &T{10})
	assertValue(t, container1, "b", &T{20})
	assertValue(t, container1, "c", &T{30})
	assertValue(t, container1, "d", &T{75})
	assertValue(t, container1, "e", &T{75})

	assertValue(t, container2, "a", &T{25})
	assertValue(t, container2, "b", &T{20})
	assertValue(t, container2, "c", &T{30})
	assertValue(t, container2, "d", &T{50}) // overlay takes precedence
	assertValue(t, container2, "e", &T{75})
}

func TestContainerWithValuesSetDuplicate(t *testing.T) {
	container, err := New().WithValues(map[interface{}]interface{}{
		"dup": struct{}{},
	})

	require.Nil(t, err)
	assert.EqualError(t, container.Set("dup", struct{}{}), `duplicate service key "dup"`)
}

type testKey1 struct{ name string }
type testKey2 struct{ name string }
type testKey3 struct{ name string }

func (k testKey1) Tag() string { return k.name }
func (k testKey2) Tag() string { return k.name }

func assertValue(t *testing.T, container *Container, key, expected interface{}) {
	value, err := container.Get(key)
	require.Nil(t, err)
	assert.Equal(t, expected, value)
}
