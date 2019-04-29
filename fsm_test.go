package gofsm_test

import (
	"context"
	"errors"
	"github.com/threeq/gofaker"
	"github.com/threeq/gofsm"
	"reflect"
	"strconv"
	"testing"
)

type CustomProcessor struct {
}

func (CustomProcessor) OnExit(ctx context.Context, state gofsm.State, event gofsm.Event) error {
	return nil
}

func (CustomProcessor) OnActionFailure(ctx context.Context, from gofsm.State, event gofsm.Event, to []gofsm.State, err error) error {
	return nil
}

func (CustomProcessor) OnEnter(ctx context.Context, state gofsm.State) error {
	return nil
}

func Test_stateMachine_Trigger(t *testing.T) {
	type fields struct {
		processor   gofsm.EventProcessor
		transitions []gofsm.Transition
		states      gofsm.StatesDef
		events      gofsm.EventsDef
	}
	type args struct {
		ctx   context.Context
		from  gofsm.State
		event gofsm.Event
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    gofsm.State
		wantErr bool
	}{
		{"Empty",
			fields{
				processor: gofsm.NoopProcessor},
			args{nil, gofsm.Start, gofsm.None},
			"", true},
		{"Not State",
			fields{
				states:    gofsm.StatesDef{"s1": "s1", "s2": "s2",},
				processor: gofsm.NoopProcessor},
			args{nil, "s3", gofsm.None},
			"", true},
		{"Not Event",
			fields{
				states:    gofsm.StatesDef{"s1": "s1", "s2": "s2",},
				events:    gofsm.EventsDef{"e1": "e1", "e2": "e2",},
				processor: gofsm.NoopProcessor},
			args{nil, "s1", "e3"},
			"", true},
		{"Not Transition",
			fields{
				states: gofsm.StatesDef{"s1": "s1", "s2": "s2", "s3": "s3"},
				events: gofsm.EventsDef{"e1": "e1", "e2": "e2",},
				transitions: []gofsm.Transition{
					{"s1", "e1", []gofsm.State{"s2"}, gofsm.NoopAction, nil},
				},
				processor: gofsm.NoopProcessor},
			args{nil, "s1", "e2"},
			"", true},

		{"Has Transition",
			fields{
				states: gofsm.StatesDef{"s1": "s1", "s2": "s2", "s3": "s3"},
				events: gofsm.EventsDef{"e1": "e1", "e2": "e2",},
				transitions: []gofsm.Transition{
					{"s1", "e1", nil, gofsm.NoopAction, nil},
				},
				processor: gofsm.NoopProcessor},
			args{nil, "s1", "e1"},
			"", false},
		{"Has Transition",
			fields{
				states: gofsm.StatesDef{"s1": "s1", "s2": "s2", "s3": "s3"},
				events: gofsm.EventsDef{"e1": "e1", "e2": "e2",},
				transitions: []gofsm.Transition{
					{"s1", "e1", []gofsm.State{"s2"}, gofsm.NoopAction, nil},
				},
				processor: gofsm.NoopProcessor},
			args{nil, "s1", "e1"},
			"s2", false},
		{"Has Transition",
			fields{
				states: gofsm.StatesDef{"s1": "s1", "s2": "s2", "s3": "s3"},
				events: gofsm.EventsDef{"e1": "e1", "e2": "e2",},
				transitions: []gofsm.Transition{
					{"s1", "e1", []gofsm.State{"s2"}, gofsm.NoopAction, nil},
					{"s1", "e2", []gofsm.State{"s3"}, gofsm.NoopAction, nil},
					{"s1", "e2", []gofsm.State{"s3"}, gofsm.NoopAction, nil},
				},
				processor: gofsm.NoopProcessor},
			args{nil, "s1", "e2"},
			"s3", false},
		{"StateMachine Processor is nil",
			fields{
				states: gofsm.StatesDef{"s1": "s1", "s2": "s2", "s3": "s3"},
				events: gofsm.EventsDef{"e1": "e1", "e2": "e2",},
				transitions: []gofsm.Transition{
					{"s1", "e1", []gofsm.State{"s2"}, gofsm.NoopAction, nil},
					{"s1", "e2", []gofsm.State{"s3"}, gofsm.NoopAction, nil},
					{"s1", "e2", []gofsm.State{"s3"}, gofsm.NoopAction, nil},
				},
				processor: nil},
			args{nil, "s1", "e2"},
			"s3", false},
		{"Transition Processor",
			fields{
				states: gofsm.StatesDef{"s1": "s1", "s2": "s2", "s3": "s3"},
				events: gofsm.EventsDef{"e1": "e1", "e2": "e2",},
				transitions: []gofsm.Transition{
					{"s1", "e1", []gofsm.State{"s2"}, gofsm.NoopAction, nil},
					{"s1", "e2", []gofsm.State{"s3"}, gofsm.NoopAction, &CustomProcessor{}},
					{"s1", "e2", []gofsm.State{"s3"}, gofsm.NoopAction, nil},
				},
				processor: gofsm.NoopProcessor},
			args{nil, "s1", "e2"},
			"s3", false},
		{"Action Error Default Processor",
			fields{
				states: gofsm.StatesDef{"s1": "s1", "s2": "s2", "s3": "s3"},
				events: gofsm.EventsDef{"e1": "e1", "e2": "e2",},
				transitions: []gofsm.Transition{
					{"s1", "e1", []gofsm.State{"s2"}, gofsm.NoopAction, nil},
					{"s1", "e2", []gofsm.State{"s3"}, func(ctx context.Context, from gofsm.State, event gofsm.Event, to []gofsm.State) (state gofsm.State, e error) {
						return "", errors.New("action error")
					}, nil},
					{"s1", "e2", []gofsm.State{"s3"}, gofsm.NoopAction, nil},
				},
				processor: gofsm.NoopProcessor},
			args{nil, "s1", "e2"},
			"", true},
		{"Action Error Customer Processor",
			fields{
				states: gofsm.StatesDef{"s1": "s1", "s2": "s2", "s3": "s3"},
				events: gofsm.EventsDef{"e1": "e1", "e2": "e2",},
				transitions: []gofsm.Transition{
					{"s1", "e1", []gofsm.State{"s2"}, gofsm.NoopAction, nil},
					{"s1", "e2", []gofsm.State{"s3"}, func(ctx context.Context, from gofsm.State, event gofsm.Event, to []gofsm.State) (state gofsm.State, e error) {
						return "", errors.New("action error")
					}, &CustomProcessor{}},
					{"s1", "e2", []gofsm.State{"s3"}, gofsm.NoopAction, nil},
				},
				processor: gofsm.NoopProcessor},
			args{nil, "s1", "e2"},
			"", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := gofsm.New("").
				States(tt.fields.states).
				Events(tt.fields.events).
				Transitions(tt.fields.transitions...).
				Processor(tt.fields.processor)
			got, err := sm.Trigger(tt.args.ctx, tt.args.from, tt.args.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("stateMachine.Trigger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("stateMachine.Trigger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkStateMachine_Trigger(b *testing.B) {
	b.StopTimer() //调用该函数停止压力测试的时间计数

	states := gofsm.StatesDef{}
	for i := 0; i < 100; i++ {
		states["s"+strconv.Itoa(i)] = "ss " + strconv.Itoa(i)
	}
	events := gofsm.EventsDef{}
	for i := 0; i < 100; i++ {
		events["e"+strconv.Itoa(i)] = "ee " + strconv.Itoa(i)
	}

	sm := gofsm.New("").
		States(states).
		Events(events).
		Processor(gofsm.NoopProcessor)

	for n := 0; n < 200; n++ {
		from := gofaker.NaturalN(0, 100)
		to := gofaker.NaturalN(0, 100)
		sm.Transitions(gofsm.Transition{
			From:   "s" + strconv.Itoa(from),
			Event:  "e" + strconv.Itoa(n),
			To:     []gofsm.State{"s" + strconv.Itoa(to)},
			Action: gofsm.NoopAction})
	}

	println(sm.Show())

	b.StartTimer() //重新开始时间
	for i := 0; i < b.N; i++ {
		s := strconv.Itoa(i % 30)
		_, _ = sm.Trigger(context.TODO(), "s"+s, "e"+s)
	}
}
