package gosm

//func TestStateMachine(t *testing.T) {
//	builder := gosm.NewBuilder()
//	builder.Trans().
//		From("").
//		To("").
//		On("").
//		When("Any", gosm.Any).
//		Action("", func(ctx context.Context, entity gosm.Entity, from, to gosm.state) error {
//			return nil
//		})
//	builder.Build("machineID")
//	machine := gosm.Get("machineID")
//
//	event := ""
//	entity := NewTestEntity("t1", "s1")
//	machine.Trigger(context.Background(), entity, event)
//}
