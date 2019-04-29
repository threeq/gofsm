package gofsm_test

import (
	"context"
	"errors"
	"fmt"
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

type OrderEventProcessor struct{}

func (*OrderEventProcessor) OnExit(ctx context.Context, state gofsm.State, event gofsm.Event) error {
	println(fmt.Sprintf("OnExit: [%v] Exit [%v] on event [%v]", ctx.Value("data"), state, event))
	return nil
}

func (*OrderEventProcessor) OnActionFailure(ctx context.Context, from gofsm.State, event gofsm.Event, to []gofsm.State, err error) error {
	println(fmt.Sprintf("OnActionFailure: [%v] do action error %v --%v--> %v", ctx.Value("data"), from, event, to))
	return nil
}

func (*OrderEventProcessor) OnEnter(ctx context.Context, state gofsm.State) error {
	println(fmt.Sprintf("OnEnter: [%v] Enter [%v]", ctx.Value("data"), state))
	return nil
}

func TestStateMachine_Example_Order(t *testing.T) {
	// 订单状态定义
	const (
		Start       = "Start"
		WaitPay     = "WaitPay"
		Paying      = "Paying"
		WaitSend    = "WaitSend"
		Sending     = "Sending"
		WaitConfirm = "WaitConfirm"
		Received    = "Received"
		PayFailure  = "PayFailure"
		Canceled    = "Canceled"
	)

	// 订单时间定义
	const (
		CreateEvent          = "Create"
		PayEvent             = "Pay"
		PaySuccessEvent      = "PaySuccess"
		PayFailureEvent      = "PayFailure"
		SendStartEvent       = "SendStart"
		SendEndEvent         = "SendEnd"
		ConfirmReceivedEvent = "SendConfirm"
		CancelEvent          = "Cancel"
	)

	doAction := func(ctx context.Context, from gofsm.State, event gofsm.Event, to []gofsm.State) (state gofsm.State, e error) {
		println(fmt.Sprintf("doAction: [%v] --%s--> %v", ctx.Value("data"), event, to))
		return to[0], nil
	}

	orderStateMachine := gofsm.New("myStateMachine").
		States(gofsm.StatesDef{
			Start:       "开始",
			WaitPay:     "待支付",
			Paying:      "支付中",
			WaitSend:    "待发货",
			Sending:     "运输中",
			WaitConfirm: "已收货",
			Received:    "已收货",
			PayFailure:  "支付失败",
			Canceled:    "已取消",
		}).
		Start([]gofsm.State{Start}).
		End([]gofsm.State{Received, Canceled}).
		Events(gofsm.EventsDef{
			CreateEvent:          "创建订单",
			PayEvent:             "支付",
			PaySuccessEvent:      "支付成功",
			PayFailureEvent:      "支付失败",
			SendStartEvent:       "发货",
			SendEndEvent:         "送达",
			ConfirmReceivedEvent: "确认收货",
			CancelEvent:          "去掉订单",
		}).
		Transitions([]gofsm.Transition{
			{Start, CreateEvent, []gofsm.State{WaitPay}, doAction, nil},
			{WaitPay, PayEvent, []gofsm.State{Paying}, doAction, nil},
			{WaitPay, CancelEvent, []gofsm.State{Canceled}, doAction, nil},
			{Paying, PaySuccessEvent, []gofsm.State{WaitSend}, doAction, nil},
			{Paying, PayFailureEvent, []gofsm.State{PayFailure}, doAction, nil},
			{PayFailure, PayEvent, []gofsm.State{Paying}, doAction, nil},
			{PayFailure, CancelEvent, []gofsm.State{Canceled}, doAction, nil},
			{WaitSend, SendStartEvent, []gofsm.State{Sending}, doAction, nil},
			{Sending, SendEndEvent, []gofsm.State{WaitConfirm}, doAction, nil},
			{WaitConfirm, ConfirmReceivedEvent, []gofsm.State{Received}, doAction, nil},
		}...)

	println(orderStateMachine.Show())

	orderStateMachine.Processor(&OrderEventProcessor{})

	order := context.WithValue(context.TODO(), "data", "order object data")
	state, err := orderStateMachine.Trigger(order, Start, CreateEvent)
	println(fmt.Sprintf("====: %v : %v", state, err))

	state, err = orderStateMachine.Trigger(order, Start, PayEvent)
	println(fmt.Sprintf("====: %v : %v", state, err))

	state, err = orderStateMachine.Trigger(order, Paying, PayFailureEvent)
	println(fmt.Sprintf("====: %v : %v", state, err))
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
