package gosm_test

import (
	"context"
	"github.com/threeq/gosm"
	"testing"
)


func TestStateMachine(t *testing.T) {
	builder := gosm.NewBuilder()
	builder.Transition().
		From("").
		To("").
		On("").
		When("Any", gosm.Any).
		Action(func(ctx context.Context, entity gosm.Entity, from, to gosm.State) error {
			return nil
		})
	builder.Build("machineID")
	machine := gosm.Get("machineID")

	event := ""
	entity := NewTestEntity("t1", "s1")
	machine.Trigger(context.Background(), entity, event)
}
