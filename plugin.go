package gosm

import (
	"context"
	"sync"
)

var (
	NoopFilter    = new(noopFilter)
)

type Event int

type LockerFactory interface {
	New(id string) sync.Locker
}

//
//
//type StateMachine struct {
//	Name          string
//	transitions   map[interface{}]map[Event][]*Transition
//	lockerFactory LockerFactory
//	filter        Filter
//
//	states []state
//}
//
//func (o *StateMachine) Transition(transitions ...*Transition) {
//	for _, trans := range transitions {
//		stateEvents, stateExist := o.transitions[trans.From.ID()]
//		if !stateExist {
//			stateEvents = make(map[Event][]*Transition)
//		}
//		transitions, eventExist := stateEvents[trans.Event]
//		if !eventExist {
//			transitions = []*Transition{}
//		}
//		transitions = append(transitions, trans)
//		stateEvents[trans.Event] = transitions
//		o.transitions[trans.From.ID()] = stateEvents
//	}
//}
//
//func (o *StateMachine) Trigger(c context.Context, entity Entity, event Event) error {
//	state := entity.State()
//	stateEvents, stateExist := o.transitions[state.ID()]
//	if !stateExist {
//		return errors.New(fmt.Sprintf("%s 状态没有定义", state))
//	}
//	transitions, eventExist := stateEvents[event]
//	if !eventExist {
//		return errors.New(fmt.Sprintf("%s - %v 没有定义", state, event))
//	}
//
//	// 支持并发控制
//	if o.lockerFactory != nil {
//		locker := o.lockerFactory.New(entity.ID())
//		locker.Lock()
//		defer locker.Unlock()
//	}
//
//	entity = o.filter.Before(c, entity, event)
//
//	var transition *Transition
//	for _, trans := range transitions {
//		if trans.Cond(c, entity, trans.From, trans.To) {
//			transition = trans
//			break
//		}
//		log.Printf("%s：条件检查失败", trans.Text())
//	}
//	if transition == nil {
//
//		return errors.New(fmt.Sprintf("%s 所有事件检查均失败", state))
//	}
//	err := transition.Action(c, entity, transition.From, transition.To)
//
//	o.filter.After(c, entity, transition, err)
//	return err
//}
//
////Show
//// Text Graph:
////      F(state = Condition1) -> NewState1 : Action1;
////      F(state = Condition2) -> NewState2 : Action2;
//// Image Graph:
////      PlantUML Code
//func (o *StateMachine) Show() string {
//	buffer := &strings.Builder{}
//	for _, tt := range o.transitions {
//		for _, transitions := range tt {
//			for _, t2 := range transitions {
//				buffer.WriteString(t2.Text())
//				buffer.WriteString("；\n")
//			}
//		}
//	}
//
//	buffer.WriteString("\n\n\n")
//	buffer.WriteString("PlantUML Code:\n\n")
//
//	// 头部信息
//	title := ""
//	if o.Name != "" {
//		title = "<b>[" + o.Name + "]</b> "
//	}
//
//	buffer.WriteString(`
//	@startuml
//	skinparam state {
//		BackgroundColor<<NFA>> Red
//	}
//	state "<font color=red><b><<DFA>></b></font>\n` + title + `state Graph" as rootGraph {
//
//	`)
//
//	// 状态的定义
//	for _, state := range o.states {
//		buffer.WriteString(fmt.Sprintf(`state "%v" as %v `, state.ID(), state.ID()))
//		buffer.WriteString("\n")
//	}
//
//	// 处理中间状态转换
//	for _, tt := range o.transitions {
//		for _, transitions := range tt {
//			for _, t2 := range transitions {
//
//				buffer.WriteString(fmt.Sprintf("%s --> %s :%v<%s>\n",
//					t2.From,
//					t2.To,
//					t2.Event,
//					t2.CondDesc))
//			}
//		}
//	}
//	buffer.WriteString("}\n")
//	buffer.WriteString("@enduml")
//
//	buffer.WriteString("\n\n\n 使用 http://www.plantuml.com/plantuml/uml/SyfFKj2rKt3CoKnELR1Io4ZDoSa70000 查看对应状态图")
//
//	return buffer.String()
//}
//
////-------------------------------------------------------------
//
//
//func Get(s string) *StateMachine {
//	return stateMachines[s]
//}
//
//func NewStateMachine(options ...Option) *StateMachine {
//	sm := &StateMachine{
//		transitions: make(map[interface{}]map[Event][]*Transition),
//		filter:      NoopFilter,
//	}
//	for _, option := range options {
//		option(sm)
//	}
//	return sm
//}

//-------------------------------------------------------------

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
