package gosm

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
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

	Name           string
	linkedMachines map[string]*StateMachine
	steps          map[*StateMachine]bool
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
	//if _, has := o.states[to.State().ID()]; !has {
	//	o.states[to.State().ID()] = to.State()
	//}

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
		return errors.New(fmt.Sprintf("%s - %v 事件没有定义", state, event))
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

func (o *StateMachine) Entry(s1 IState, desc string, action ...Action) StateEntry {
	s1.Bind(o)
	o.starts = append(o.starts, s1)
	return s1.Entry(desc, action...)
}

func (o *StateMachine) Exit(s1 IState, event Event, desc string, condition ...Condition) *StateExit {
	s1.Bind(o)
	o.ends = append(o.ends, s1)
	return s1.Exit(event, desc, condition...)
}

func (o *StateMachine) Fork(exit *StateExit) *ForkStateExit {
	return &ForkStateExit{exit, o}
}

func (o *StateMachine) end(actions ...Action) StateEntry {
	return &normalStateEntry{
		state:   o.State("[*]"),
		actions: actions,
		desc:    "end",
	}
}

//Graph
// TODO 支持嵌套
// choice、fork 显示
// Text Graph:
//      F(state = Condition1) -> NewState1 : Action1;
//      F(state = Condition2) -> NewState2 : Action2;
// Image Graph:
//      PlantUML Code
func (o *StateMachine) Graph(steps map[*StateMachine]bool) (string, string) {

	// 处理循环依赖
	if _, ok := steps[o]; ok {
		return "", ""
	}
	steps[o] = true
	o.steps = steps

	statesBuffer := &strings.Builder{}
	// 头部信息
	title := ""
	if o.Name != "" {
		title = "<b>[" + o.Name + "]</b> "
	}

	statesBuffer.WriteString(`state "<font color=red><b><<DFA>></b></font>\n` +
		title + `state Graph" as ` + o.Name + " {\n")

	// 状态的定义
	var machineStates []string
	// 处理中间状态转换
	var stateLines []string
	var transLines []string
	for stateID, tt := range o.transitions {
		for event, transitions := range tt {

			if len(transitions) == 1 {
				// simple transition
				t2 := transitions[0]
				m, line := t2.Graph()
				if m != "" {
					machineStates = append(machineStates, m)
				}
				transLines = append(transLines, line)
			} else if len(transitions) > 1 {
				// choice
				choiceID := fmt.Sprintf("%v_%v_choice", stateID, event)
				choice := State(choiceID, "choice")
				choice.Bind(o)

				t1 := transitions[0]
				otherMachineStates, trans := choice.Entry("").Graph(t1.exit)
				if otherMachineStates != "" {
					machineStates = append(machineStates, otherMachineStates)
				}
				transLines = append(transLines, trans)

				choiceExit := choice.Exit(t1.exit.event, t1.exit.desc)
				for _, t2 := range transitions {
					m, line := t2.Graph(choiceExit)
					if m != "" {
						machineStates = append(machineStates, m)
					}
					transLines = append(transLines, line)
				}
			}

		}
	}

	for _, ss := range o.states {
		if ss.ID() == "[*]" {
			continue
		}

		if s, ok := ss.(*state); ok && s.stereotype != "" {
			stateLines = append(stateLines, fmt.Sprintf(`    state "%v" as %v <<%s>>`, ss.ID(), ss.ID(), s.stereotype))
			continue
		}

		stateLines = append(stateLines, fmt.Sprintf(`    state "%v" as %v`, ss.ID(), ss.ID()))
	}
	statesBuffer.WriteString(strings.Join(stateLines, "\n"))
	statesBuffer.WriteString("\n}\n")
	statesBuffer.WriteString(strings.Join(machineStates, "\n"))

	return statesBuffer.String(), strings.Join(transLines, "\n")
}

func (o *StateMachine) Show() {
	steps := make(map[*StateMachine]bool)
	state, trans := o.Graph(steps)
	println("\n", state, "\n", trans, "\n")
}

//---------------------------------------------------------------------------------
func NewMachine(name string, options ...Option) *StateMachine {
	sm := &StateMachine{
		Name:        name,
		transitions: make(map[interface{}]map[Event][]*ConditionLinker),
		states:      make(map[interface{}]IState),
		filter:      NoopFilter,
	}
	for _, option := range options {
		option(sm)
	}
	stateMachines[sm.Name] = sm
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
