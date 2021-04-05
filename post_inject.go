package service

import "context"

// PostInject is a marker interface for injectable objects which should
// perform some action after injection of services.
type PostInject interface {
	PostInject(ctx context.Context) error
}
