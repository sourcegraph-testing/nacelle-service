package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostInject(t *testing.T) {
	container := New()
	obj := &testPostInjectProcess{}
	container.Set("value", &TI{42})
	err := Inject(context.Background(), container, obj)
	require.Nil(t, err)
	assert.Equal(t, 42.0, obj.FValue.val)
}

type TI struct{ val int }
type TF struct{ val float64 }

type testPostInjectProcess struct {
	IValue *TI `service:"value"`
	FValue *TF
}

var _ PostInject = &testPostInjectProcess{}

func (p *testPostInjectProcess) PostInject(ctx context.Context) error {
	p.FValue = &TF{float64(p.IValue.val)}
	return nil
}

func TestPostInjectChain(t *testing.T) {
	container := New()
	obj := &testPostInjectProcessParent{}
	process := &testPostInjectProcess{}

	container.Set("value", &TI{42})
	container.Set("process", process)
	container.Set("services", container)

	err := Inject(context.Background(), container, obj)
	require.Nil(t, err)
	assert.Equal(t, 42.0, process.FValue.val)
}

type testPostInjectProcessParent struct {
	Services *Container             `service:"services"`
	Child    *testPostInjectProcess `service:"process"`
}

var _ PostInject = &testPostInjectProcessParent{}

func (p *testPostInjectProcessParent) PostInject(ctx context.Context) error {
	return Inject(ctx, p.Services, p.Child)
}

func TestPostInjectError(t *testing.T) {
	container := New()
	obj := &testPostInjectProcessError{}
	err := Inject(context.Background(), container, obj)
	assert.EqualError(t, err, "oops")
}

type testPostInjectProcessError struct{}

var _ PostInject = &testPostInjectProcessError{}

func (p *testPostInjectProcessError) PostInject(ctx context.Context) error {
	return fmt.Errorf("oops")
}
