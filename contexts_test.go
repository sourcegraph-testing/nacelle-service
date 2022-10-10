package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithAndFromContext(t *testing.T) {
	type T1 struct{ val int }

	t.Run("from context without value", func(t *testing.T) {
		ctx := context.Background()
		assert.Nil(t, FromContext(ctx))
	})
	t.Run("from context that has value", func(t *testing.T) {
		container := New()
		container.Set("a", &T1{10})

		ctx := context.Background()
		ctx = WithContainer(ctx, container)

		ctxContainer := FromContext(ctx)
		require.NotNil(t, ctxContainer)

		assertValue(t, ctxContainer, "a", &T1{10})
	})
}
