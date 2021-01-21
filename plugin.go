package gosm

import (
	"context"
	"sync"
)

var (
	NoopFilter    = new(noopFilter)
)

type LockerFactory interface {
	New(id string) sync.Locker
}

// Filter Aspect
type Filter interface {
	Before(ctx context.Context, entity Entity, event Event) Entity
	After(ctx context.Context, entity Entity, trans *Transition, result error)
}

type noopFilter struct {
}

func (n *noopFilter) Before(_ context.Context, entity Entity, _ Event) Entity {
	return entity
}

func (n *noopFilter) After(_ context.Context, _ Entity, _ *Transition, _ error) {

}
