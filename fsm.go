package gosm

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
)

var (
	stateMachines = make(map[string]*StateMachine)
	Any           = func(ctx context.Context, entity Entity, from, to State) bool { return true }
	Noop          = func(ctx context.Context, entity Entity, from, to State) error { return nil }
	NoopFilter    = new(noopFilter)
)

type State = string
type Event = string
type Entity interface {
	ID() string
	State() State
}
type Condition = func(ctx context.Context, entity Entity, from, to State) bool
type Action = func(ctx context.Context, entity Entity, from, to State) error

type LockerFactory interface {
	New(id string) sync.Locker
}

type Transition struct {
	From       State
	Event      Event
	To         State
	Cond       Condition
	Action     Action
	CondDesc   string
	ActionDesc string
}

func (owner *Transition) Text() string {
	return fmt.Sprintf("F(%s = %s<%s>) -> %s: %s",
		owner.From, owner.Event, owner.CondDesc, owner.To, owner.ActionDesc)
}

// Filter Aspect
type Filter interface {
	Before(ctx context.Context, entity Entity, event Event) Entity
	After(ctx context.Context, entity Entity, trans *Transition, result error)
}

type StateMachine struct {
	Name          string
	transitions   map[State]map[Event][]*Transition
	lockerFactory LockerFactory
	filter        Filter

	states []State
}

func (owner *StateMachine) Transition(transitions ...*Transition) {
	for _, trans := range transitions {
		stateEvents, stateExist := owner.transitions[trans.From]
		if !stateExist {
			stateEvents = make(map[Event][]*Transition)
		}
		transitions, eventExist := stateEvents[trans.Event]
		if !eventExist {
			transitions = []*Transition{}
		}
		transitions = append(transitions, trans)
		stateEvents[trans.Event] = transitions
		owner.transitions[trans.From] = stateEvents
	}
}

func (owner *StateMachine) Trigger(c context.Context, entity Entity, event Event) error {
	state := entity.State()
	stateEvents, stateExist := owner.transitions[state]
	if !stateExist {
		return errors.New(fmt.Sprintf("%s 状态没有定义", state))
	}
	transitions, eventExist := stateEvents[event]
	if !eventExist {
		return errors.New(fmt.Sprintf("%s - %s 没有定义", state, event))
	}

	// 支持并发控制
	if owner.lockerFactory != nil {
		locker := owner.lockerFactory.New(entity.ID())
		locker.Lock()
		defer locker.Unlock()
	}

	entity = owner.filter.Before(c, entity, event)

	var transition *Transition
	for _, trans := range transitions {
		if trans.Cond(c, entity, trans.From, trans.To) {
			transition = trans
			break
		}
		log.Printf("%s：条件检查失败", trans.Text())
	}
	if transition == nil {

		return errors.New(fmt.Sprintf("%s 所有事件检查均失败", state))
	}
	err := transition.Action(c, entity, transition.From, transition.To)

	owner.filter.After(c, entity, transition, err)
	return err
}

//Show
// Text Graph:
//      F(State = Condition1) -> NewState1 : Action1;
//      F(State = Condition2) -> NewState2 : Action2;
// Image Graph:
//      PlantUML Code
func (owner *StateMachine) Show() string {
	buffer := &strings.Builder{}
	for _, tt := range owner.transitions {
		for _, transitions := range tt {
			for _, t2 := range transitions {
				buffer.WriteString(t2.Text())
				buffer.WriteString("；\n")
			}
		}
	}

	buffer.WriteString("\n\n\n")
	buffer.WriteString("PlantUML Code:\n\n")

	// 头部信息
	title := ""
	if owner.Name != "" {
		title = "<b>[" + owner.Name + "]</b> "
	}

	buffer.WriteString(`
	@startuml
	skinparam state {
		BackgroundColor<<NFA>> Red
	}
	State "<font color=red><b><<DFA>></b></font>\n` + title + `State Graph" as rootGraph {

	`)

	// 状态的定义
	for _, state := range owner.states {
		buffer.WriteString(fmt.Sprintf(`state "%s" as %s `, state, state))
		buffer.WriteString("\n")
	}

	// 处理中间状态转换
	for _, tt := range owner.transitions {
		for _, transitions := range tt {
			for _, t2 := range transitions {

				buffer.WriteString(fmt.Sprintf("%s --> %s :%s<%s>\n",
					t2.From,
					t2.To,
					t2.Event,
					t2.CondDesc))
			}
		}
	}
	buffer.WriteString("}\n")
	buffer.WriteString("@enduml")

	buffer.WriteString("\n\n\n 使用 http://www.plantuml.com/plantuml/uml/SyfFKj2rKt3CoKnELR1Io4ZDoSa70000 查看对应状态图")

	return buffer.String()
}

//-------------------------------------------------------------

type Option func(*StateMachine)

func Get(s string) *StateMachine {
	return stateMachines[s]
}

func NewStateMachine(options ...Option) *StateMachine {
	sm := &StateMachine{
		transitions: make(map[State]map[Event][]*Transition),
		filter:      NoopFilter,
	}
	for _, option := range options {
		option(sm)
	}
	return sm
}

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

//-------------------------------------------------------------

type noopFilter struct {
}

func (n *noopFilter) Before(_ context.Context, entity Entity, _ Event) Entity {
	return entity
}

func (n *noopFilter) After(_ context.Context, _ Entity, _ *Transition, _ error) {

}
