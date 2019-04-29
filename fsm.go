package gofsm

import (
	"context"
	"fmt"
	"qiniupkg.com/x/errors.v7"
	"sort"
	"strings"
)

type State = string
type Event = string
type StatesDef map[State]string
type EventsDef map[Event]string
type Action func(ctx context.Context, from State, event Event, to []State) (State, error)
type EventProcessor interface {
	OnExit(ctx context.Context, state State, event Event) error
	OnActionFailure(ctx context.Context, from State, event Event, to []State, err error) error
	OnEnter(ctx context.Context, state State) error
}
type Transition struct {
	From      State
	Event     Event
	To        []State
	Action    Action
	Processor EventProcessor
}

/**
状态机执行表述图
有限状态机
	- 确定状态机
	- 非确定状态机
*/
type stateGraph struct {
	name        string // 状态图名称
	start       []State
	end         []State
	states      StatesDef
	events      EventsDef
	transitions map[State]map[Event]*Transition
}

/**
状态机
*/
type stateMachine struct {
	processor EventProcessor
	sg        *stateGraph
}

/**
默认实现
*/
type DefaultProcessor struct{}

func (*DefaultProcessor) OnExit(ctx context.Context, state State, event Event) error {
	//log.Printf("exit [%s]", state)
	return nil
}

func (*DefaultProcessor) OnActionFailure(ctx context.Context, from State, event Event, to []State, err error) error {
	//log.Printf("failure %s -(%s)-> [%s]: (%s)", from, event, strings.Join(to, "|"), err.Error())
	return nil
}

func (*DefaultProcessor) OnEnter(ctx context.Context, state State) error {
	//log.Printf("enter [%s]", state)
	return nil
}

/**
默认值定义
*/
const Start = "[*]"
const End = "[*]"
const None = ""

var NoopAction Action = func(ctx context.Context, from State, event Event, to []State) (State, error) {
	if to == nil || len(to) == 0 {
		return None, nil
	}
	return to[0], nil
}
var NoopProcessor = &DefaultProcessor{}

/**
创建一个状态机执行器
*/
func New(name string) *stateMachine {
	return (&stateMachine{
		sg: &stateGraph{
			transitions: map[State]map[Event]*Transition{},
		}}).Name(name)
}

/**
设置所有状态
*/
func (sm *stateMachine) States(states StatesDef) *stateMachine {
	sm.sg.states = states
	return sm
}

/**
设置所有时间
*/
func (sm *stateMachine) Events(events EventsDef) *stateMachine {
	sm.sg.events = events
	return sm
}

func (sm *stateMachine) Name(s string) *stateMachine {
	sm.sg.name = s
	return sm
}

func (sm *stateMachine) Start(start []State) *stateMachine {
	sm.sg.start = start
	return sm
}

func (sm *stateMachine) End(end []State) *stateMachine {
	sm.sg.end = end
	return sm
}

func (sm *stateMachine) Processor(processor EventProcessor) *stateMachine {
	sm.processor = processor
	return sm
}

/**
添加状态转换
*/
func (sm *stateMachine) Transitions(transitions ...Transition) *stateMachine {
	for index := range transitions {
		newTransfer := &transitions[index]
		events, ok := sm.sg.transitions[newTransfer.From]
		if !ok {
			events = map[Event]*Transition{}
			sm.sg.transitions[newTransfer.From] = events
		}
		if transfer, ok := events[newTransfer.Event]; ok {
			transfer.To = append(transfer.To, newTransfer.To...)
			// 去掉重复
			sort.Strings(transfer.To)
			transfer.To = removeDuplicatesAndEmpty(transfer.To)
			events[newTransfer.Event] = transfer
		} else {
			events[newTransfer.Event] = newTransfer
		}
	}
	return sm
}

func removeDuplicatesAndEmpty(a []State) (ret []State) {
	aLen := len(a)
	for i := 0; i < aLen; i++ {
		if (i > 0 && a[i-1] == a[i]) || len(a[i]) == 0 {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}

/**
触发状态转换
*/
func (sm *stateMachine) Trigger(ctx context.Context, from State, event Event) (State, error) {
	if _, ok := sm.sg.states[from]; !ok {
		return "", errors.New("状态机不包含状态" + from)
	}
	if _, ok := sm.sg.events[event]; !ok {
		return "", errors.New("状态机不包含事件 " + event)
	}
	if transfer, ok := sm.sg.transitions[from][event]; ok {

		processor := sm.processor
		// 离开状态处理，转换之前
		if transfer.Processor != nil {
			processor = transfer.Processor
		}
		if processor == nil {
			processor = NoopProcessor
		}

		_ = processor.OnExit(ctx, from, event)

		to, err := transfer.Action(ctx, from, event, transfer.To)
		if err != nil {
			// 转换执行错误处理
			_ = processor.OnActionFailure(ctx, from, event, transfer.To, err)
			return to, err
		}
		// 进入状态处理，转换之后
		_ = processor.OnEnter(ctx, to)

		return to, err
	}
	return "", errors.New(fmt.Sprintf("没有定义状态转换事件 [%v --%v--> ???]", from, event))

}

/**
输出图的显示内容
输出 PlantUML 显示 URL
*/
func (sm *stateMachine) Show() string {
	return sm.sg.show()
}

func (transfer Transition) String() string {
	return fmt.Sprintf("%s --> %s: %s", transfer.From, transfer.To, transfer.Event)
}

/**
输出图的显示内容
输出 PlantUML 显示 URL
*/
func (sg *stateGraph) show() string {
	// 头部信息
	title := ""
	if sg.name != "" {
		title = "[" + sg.name + "] "
	}
	// 状态的定义
	var stateLines []string
	for state, desc := range sg.states {
		stateLine := string(state)
		if desc != "" {
			stateLine = fmt.Sprintf("%s: %s", state, desc)
		}

		stateLines = append(stateLines, stateLine)
	}
	statesDef := strings.Join(stateLines, "\n")

	// 状态转换描述
	var transferLines []string
	// 开始状态处理
	if sg.start != nil && len(sg.start) > 0 {
		for _, event := range sg.start {
			transferLines = append(transferLines,
				fmt.Sprintf("%s --> %s",
					Start,
					event))
		}
	}
	// 处理中间状态转换
	for from, events := range sg.transitions {
		for event, transfer := range events {
			if event != "" {
				desc := sg.events[event]
				event = "(" + string(event) + ") "
				if desc != "" {
					event += desc
				}
			}
			// plantUml 格式
			if event != "" {
				event = ": " + event
			}

			for j := 0; j < len(transfer.To); j++ {
				to := transfer.To[j]
				transferLines = append(transferLines,
					fmt.Sprintf("%s --> %s %s",
						from,
						to,
						event))
			}
		}
	}
	// 结束状态处理
	if sg.end != nil && len(sg.end) > 0 {
		for _, event := range sg.end {
			transferLines = append(transferLines,
				fmt.Sprintf("%s --> %s",
					event,
					End))
		}
	}
	transitionsDef := strings.Join(transferLines, "\n")

	// 生成 plantUml script
	raw := `
	@startuml
	
	State "` + title + `State Graph" as rootGraph {
		%s

		%s
	}

	@enduml
	`
	raw = fmt.Sprintf(raw, statesDef, transitionsDef)

	// 输出 plantUml 和 在线生成图标地址
	plantText := encode(raw)
	imgUrl := "https://www.plantuml.com/plantuml/img/" + plantText
	svgUrl := "https://www.plantuml.com/plantuml/svg/" + plantText
	format := "\nPlantUml Script:\n%s\n\nOnline Graph:\n\tImg: %s\n\tSvg: %s"
	return fmt.Sprintf(format, raw, imgUrl, svgUrl)
}
