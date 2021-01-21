package gosm

import (
	"context"
	"fmt"
)

var (
	Any  = func(ctx context.Context, entity Entity, from, to IState) bool { return true }
	Noop = func(ctx context.Context, entity Entity, from, to IState) error { return nil }
)

type Condition func(ctx context.Context, entity Entity, from, to IState) bool
type Action func(ctx context.Context, entity Entity, from, to IState) error

type Event interface{}

type IState interface {
	ID() interface{}
	Entry(desc string, actions ...Action) StateEntry
	Exit(event Event, desc string, condition ...Condition) *StateExit
	Bind(machine *StateMachine)
	Machine() *StateMachine
}

func State(v interface{}, p ...string) IState {
	if len(p) > 0 {
		return &state{value: v, stereotype: p[0]}
	}
	return &state{value: v}
}

type state struct {
	value      interface{}
	stereotype string
	machine    *StateMachine
}

func (o *state) ID() interface{} {
	return o.value
}

func (o *state) Bind(machine *StateMachine) {
	o.machine = machine
	machine.states[o.ID()] = o
}

func (o *state) Machine() *StateMachine {
	return o.machine
}

func (o *state) Entry(desc string, actions ...Action) StateEntry {
	return &normalStateEntry{state: o, actions: actions, desc: desc}
}

func (o *state) Exit(event Event, desc string, condition ...Condition) *StateExit {
	cond := Any
	d := "Any"
	if len(condition) > 0 {
		cond = condition[0]
		d = desc
	}

	return &StateExit{
		state: o,
		event: event,
		cond:  cond,
		desc:  d,
	}
}

type StateEntry interface {
	State() IState
	Action(ctx context.Context, entity Entity, from, to IState) error
	Desc() string
	Graph(exit *StateExit) (string, string)
}

type normalStateEntry struct {
	state   IState
	actions []Action
	desc    string
}

func (o *normalStateEntry) State() IState {
	return o.state
}

func (o *normalStateEntry) Action(ctx context.Context, entity Entity, from, to IState) error {
	for _, action := range o.actions {
		err := action(ctx, entity, from, to)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *normalStateEntry) Graph(exit *StateExit) (string, string) {
	// 同一个状态机
	if exit.state.Machine() == nil ||
		o.State().Machine() == nil ||
		exit.state.Machine() == o.State().Machine() {
		return "", fmt.Sprintf("%v --> %v :%v<%v>",
			exit.state.ID(),
			o.state.ID(),
			exit.event,
			exit.desc)
	}
	transLine := fmt.Sprintf("%v --> %v :%v<%v>",
		exit.state.ID(),
		o.state.ID(),
		exit.event,
		exit.desc)

	m, l := o.State().Machine().Graph(exit.state.Machine().steps)
	return m, transLine + "\n" + l

}

func (o *normalStateEntry) Desc() string {
	return o.desc
}

type StateExit struct {
	state     IState
	event     Event
	cond      Condition
	desc      string
	nextEntry StateEntry
}

func (o *StateExit) End(actions ...Action) {
	o.state.Machine().Trans(o, o.state.Machine().end(actions...))
}

func (o *StateExit) Link(entry StateEntry) {
	o.nextEntry = entry
	o.state.Machine().Trans(o, entry)
}

type ConditionLinker struct {
	exit  *StateExit
	entry StateEntry
}

func (o *ConditionLinker) Text() string {
	return fmt.Sprintf("F(%v = %v<%v>) -> %v",
		o.exit.state.ID(), o.exit.event, o.exit.desc, o.entry.State().ID())
}

func (o *ConditionLinker) Graph(exit ...*StateExit) (string, string) {
	e := o.exit
	if len(exit) > 0 {
		e = exit[0]
	}

	return o.entry.Graph(e)
}

type Entity interface {
	ID() string
	State() IState
}
