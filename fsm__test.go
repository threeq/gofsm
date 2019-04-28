package gofsm

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *stateMachine
	}{
		{"new", &stateMachine{processor: NoopProcessor, sg: &stateGraph{transitions: map[State]map[Event]*Transition{}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(""); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New(nil) = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stateMachine_Show(t *testing.T) {
	type fields struct {
		name        string
		states      StatesDef
		events      EventsDef
		transitions []Transition
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"Empty Graph", fields{
			"",
			nil,
			nil,
			[]Transition{},
		}, `State Graph`},
		{"Sample State Graph", fields{
			"Sample",
			StatesDef{
				"s1": "stat 1",
				"s2": "stat 2",
			},
			nil,
			[]Transition{
				{Start, "Start", []State{"s1"}, NoopAction, nil},
				{Start, None, []State{"s2"}, NoopAction, nil},
				{"s2", "Execute", []State{End, "44"}, NoopAction, nil},
				{"s2", "Execute", []State{End, "33"}, NoopAction, NoopProcessor},
				{"s2", "SS", []State{End, "11", "22"}, NoopAction, nil},
			},
		}, `State "[Sample] State Graph" as rootGraph`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := New("").
				Name(tt.fields.name).
				States(tt.fields.states).
				Events(tt.fields.events).
				Transitions(tt.fields.transitions...)
			println(sm.Show())
			if got := sm.Show(); !strings.Contains(got, tt.want) {
				t.Errorf("stateMachine.Show() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stateMachine_States(t *testing.T) {
	type args struct {
		states StatesDef
	}
	tests := []struct {
		name string
		args args
		want *stateMachine
	}{
		{"Empty", args{nil}, New("")},
		{"Has States", args{StatesDef{"a": "a", "b": "b"}}, &stateMachine{
			processor: NoopProcessor,
			sg: &stateGraph{
				states:      StatesDef{"a": "a", "b": "b"},
				transitions: map[State]map[Event]*Transition{},
			}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := New("")
			if got := sm.States(tt.args.states); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("stateMachine.States() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stateMachine_Events(t *testing.T) {

	type args struct {
		events EventsDef
	}
	tests := []struct {
		name string
		args args
		want *stateMachine
	}{
		{"Empty", args{nil}, New("")},
		{"Has Events", args{EventsDef{"a": "a", "b": "b"}}, &stateMachine{
			processor: NoopProcessor,
			sg: &stateGraph{
				events:      EventsDef{"a": "a", "b": "b"},
				transitions: map[State]map[Event]*Transition{},
			}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := New("")
			if got := sm.Events(tt.args.events); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("stateMachine.Events() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stateMachine_Transitions(t *testing.T) {
	// 数据定义
	ts := []Transition{
		{Start, None, []State{End}, NoopAction, nil},
		{Start, "event1", []State{End, "test2"}, NoopAction, nil},
	}

	// table
	type args struct {
		transitions []Transition
	}
	tests := []struct {
		name string
		args args
		want *stateMachine
	}{
		{"Empty", args{nil}, New("")},
		{"Has Transitions", args{ts}, New("").Transitions(ts...)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := New("")
			sm.Transitions(tt.args.transitions...)
			println(fmt.Sprintf("%v", sm.sg.transitions))
			if got := sm.Transitions(tt.args.transitions...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("stateMachine.Transitions() = %v, want %v", got, tt.want)
			}
		})
	}
}
