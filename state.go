package gosm

import "context"

var (
	Any  = func(ctx context.Context, entity Entity, from, to IState) bool { return true }
	Noop = func(ctx context.Context, entity Entity, from, to IState) error { return nil }
)

type Condition func(ctx context.Context, entity Entity, from, to IState) bool
type Action func(ctx context.Context, entity Entity, from, to IState) error

type IState interface {
	ID() interface{}
	Entry(actions ...Action) StateEntry
	Exit(event Event, condition ...Condition) *StateExit
	Bind(machine *StateMachine)
	Machine() *StateMachine
}

func State(v interface{}) IState {
	return &state{value: v}
}

type state struct {
	value   interface{}
	machine *StateMachine
}

func (o *state) ID() interface{} {
	return o.value
}

func (o *state) Bind(machine *StateMachine) {
	o.machine = machine
}

func (o *state) Machine() *StateMachine {
	return o.machine
}

func (o *state) Entry(actions ...Action) StateEntry {
	//TODO desc 输入
	return &SimpleStateEntry{state: o, actions: actions, desc: ""}
}

func (o *state) Exit(event Event, condition ...Condition) *StateExit {
	cond := Any
	if len(condition) > 0 {
		cond = condition[0]
	}

	return &StateExit{
		state: o,
		event: event,
		cond:  cond,
		//TODO desc 输入
		desc: "",
	}
}

type SimpleStateEntry struct {
	state   IState
	actions []Action
	desc    string
}

func (o *SimpleStateEntry) State() IState {
	return o.state
}

func (o *SimpleStateEntry) Action(ctx context.Context, entity Entity, from, to IState) error {
	for _, action := range o.actions {
		err := action(ctx, entity, from, to)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *SimpleStateEntry) Desc() string {
	return o.desc
}

type StateEntry interface {
	State() IState
	Action(ctx context.Context, entity Entity, from, to IState) error
	Desc() string
}

type StateExit struct {
	state     IState
	event     Event
	cond      Condition
	desc      string
	hasLinked bool
	nextEntry StateEntry
}

func (o *StateExit) End(actions ...Action) {
	o.hasLinked = false
	o.state.Machine().Trans(o, o.state.Machine().end(actions...))
}

func (o *StateExit) Link(entry StateEntry) {
	o.hasLinked = true
	o.nextEntry = entry
	o.state.Machine().Trans(o, entry)
}

type ConditionLinker struct {
	exit  *StateExit
	entry StateEntry
}

func (o *ConditionLinker) Text() string {
	// TODO 文本显示
	return ""
}

type Entity interface {
	ID() string
	State() IState
}
