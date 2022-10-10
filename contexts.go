package service

import "context"

type containerKeyType struct{}

var containerKey = containerKeyType{}

func WithContainer(ctx context.Context, container *Container) context.Context {
	return context.WithValue(ctx, containerKey, container)
}

func FromContext(ctx context.Context) *Container {
	if v, ok := ctx.Value(containerKey).(*Container); ok {
		return v
	}
	return nil
}
