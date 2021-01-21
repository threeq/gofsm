package gosm

import (
	"context"
	"errors"
	"fmt"
	"log"
)

type Transition struct {
	From       IState
	Event      Event
	To         IState
	Cond       Condition
	Action     Action
	CondDesc   string
	ActionDesc string
}

func (o *Transition) Text() string {
	return fmt.Sprintf("F(%s = %v<%s>) -> %s: %s",
		o.From, o.Event, o.CondDesc, o.To, o.ActionDesc)
}

type StateMachine struct {
	starts      []IState
	ends        []IState
	states      map[interface{}]IState
	transitions map[interface{}]map[Event][]*ConditionLinker

	lockerFactory LockerFactory
	filter        Filter

	Name string
}

func (o *StateMachine) Trans(from *StateExit, to StateEntry) {
	stateEvents, stateExist := o.transitions[from.state.ID()]
	if !stateExist {
		stateEvents = make(map[Event][]*ConditionLinker)
	}
	transitions, eventExist := stateEvents[from.event]
	if !eventExist {
		transitions = []*ConditionLinker{}
	}
	transitions = append(transitions, &ConditionLinker{exit: from, entry: to})
	stateEvents[from.event] = transitions
	o.transitions[from.state.ID()] = stateEvents

	if _, has := o.states[from.state.ID()]; !has {
		o.states[from.state.ID()] = from.state
	}
	if _, has := o.states[to.State().ID()]; !has {
		o.states[to.State().ID()] = to.State()
	}

	from.state.Bind(o)
}

func (o *StateMachine) Trigger(c context.Context, entity Entity, event Event) error {
	state := entity.State()
	stateEvents, stateExist := o.transitions[state.ID()]
	if !stateExist {
		return errors.New(fmt.Sprintf("%s 状态没有定义", state))
	}
	transitions, eventExist := stateEvents[event]
	if !eventExist {
		return errors.New(fmt.Sprintf("%s - %v 没有定义", state, event))
	}

	// 支持并发控制
	if o.lockerFactory != nil {
		locker := o.lockerFactory.New(entity.ID())
		locker.Lock()
		defer locker.Unlock()
	}

	entity = o.filter.Before(c, entity, event)

	// 出 状态 条件判断
	var transition *ConditionLinker
	for _, trans := range transitions {
		if trans.exit.cond(c, entity, trans.exit.state, trans.entry.State()) {
			transition = trans
			break
		}
		log.Printf("%s：条件检查失败", trans.Text())
	}
	if transition == nil {
		return errors.New(fmt.Sprintf("%s 所有事件检查均失败", state))
	}

	// 进 状态 操作逻辑
	err := transition.entry.Action(c, entity, transition.exit.state, transition.entry.State())

	o.filter.After(c, entity, &Transition{
		From: transition.exit.state, Event: transition.exit.event, CondDesc: transition.exit.desc,
		To: transition.entry.State(), ActionDesc: transition.entry.Desc(),
	}, err)
	return err
}

func (o *StateMachine) State(v interface{}) IState {
	s, exist := o.states[v]
	if !exist {
		s = &state{value: v, machine: o}
		o.states[v] = s
	}
	return s
}

func (o *StateMachine) Entry(s1 IState, action ...Action) StateEntry {
	s1.Bind(o)
	o.starts = append(o.starts, s1)
	return s1.Entry(action...)
}

func (o *StateMachine) Exit(s1 IState, event Event, condition ...Condition) *StateExit {
	s1.Bind(o)
	o.ends = append(o.ends, s1)
	return s1.Exit(event, condition...)
}

func (o *StateMachine) end(actions ...Action) StateEntry {
	return &SimpleStateEntry{
		state:   o.State("[*]"),
		actions: actions,
		desc:    "end",
	}
}

//---------------------------------------------------------------------------------
func NewMachine(options ...Option) *StateMachine {
	sm := &StateMachine{
		transitions: make(map[interface{}]map[Event][]*ConditionLinker),
		states:      make(map[interface{}]IState),
		filter:      NoopFilter,
	}
	for _, option := range options {
		option(sm)
	}
	return sm
}

type Option func(*StateMachine)

func Locker(lockerFactory LockerFactory) Option {
	return func(machine *StateMachine) {
		machine.lockerFactory = lockerFactory
	}
}

func Aspect(filter Filter) Option {
	return func(machine *StateMachine) {
		machine.filter = filter
	}
}
