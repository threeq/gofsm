package gofsm_test

import (
	"context"
	"github.com/threeq/gofsm"
	"reflect"
	"testing"
)

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
			args{nil, "e3", gofsm.None},
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
