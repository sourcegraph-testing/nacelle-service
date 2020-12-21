package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOverlayContainerGet(t *testing.T) {
	type T struct {
		val int
	}

	container := NewServiceContainer()
	container.Set("a", &T{10})
	container.Set("b", &T{20})
	container.Set("c", &T{30})

	overlay := Overlay(container, map[interface{}]interface{}{
		"a": &T{40},
		"d": &T{50},
	})

	value1, err1 := overlay.Get("a")
	require.Nil(t, err1)
	assert.Equal(t, &T{40}, value1)

	value2, err2 := overlay.Get("b")
	require.Nil(t, err2)
	assert.Equal(t, &T{20}, value2)

	value3, err3 := overlay.Get("c")
	require.Nil(t, err3)
	assert.Equal(t, &T{30}, value3)

	value4, err4 := overlay.Get("d")
	require.Nil(t, err4)
	assert.Equal(t, &T{50}, value4)
}

func TestOverlayContainerInject(t *testing.T) {
	type T1 struct {
		val int
	}
	type T2 struct {
		A *T1 `service:"a"`
		B *T1 `service:"b"`
		C *T1 `service:"c"`
		D *T1 `service:"d"`
	}

	container := NewServiceContainer()
	container.Set("a", &T1{10})
	container.Set("b", &T1{20})
	container.Set("c", &T1{30})

	overlay := Overlay(container, map[interface{}]interface{}{
		"a": &T1{40},
		"d": &T1{50},
	})

	obj := &T2{}
	err := overlay.Inject(obj)
	require.Nil(t, err)
	assert.Equal(t, 40, obj.A.val)
	assert.Equal(t, 20, obj.B.val)
	assert.Equal(t, 30, obj.C.val)
	assert.Equal(t, 50, obj.D.val)
}
