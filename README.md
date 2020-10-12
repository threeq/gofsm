# gofsm

[![Build Status](https://travis-ci.org/threeq/gofsm.svg?branch=master)](https://travis-ci.org/threeq/gofsm) [![codecov](https://codecov.io/gh/threeq/gofsm/branch/master/graph/badge.svg)](https://codecov.io/gh/threeq/gofsm)

一个简单、高效的有限状态机执行框架

## 使用

例如 电商网站中的 订单状态（这里做了简化）
![](http://www.plantuml.com/plantuml/png/ZPD1QzL04CVFsKynz6I4qdk9eVgqdXJfEtYeFQp9jWrDDcMJ24M4UX2UYct5Ae87KIez-KGfLNmIlypR6B-5Exix43o8jvcTt_xdFxD9i5BLNDLDaREWsidaBbUy07DM2xZF0e0hFDdPKcKZqr6PbogARgvUZcDO4oaB7h1WRCc5QBEKDIH8N58YZQExSHHTHJ9QCk4IbkCxqXol5tlspWsUR6TIR60TdCfrnNUt5u1NeCgojXbw22hNOun6RTb60Clwnxu-VSfy_HRVo-IM1LneYExuqtpsUlxj8q6tULQFXKmjHWbAjO_quVF-x9J0DP68x9vm82K8VltI7Pzxa1HDFylvsEcvteHX7xBdOuFrFu_wziLV_aQbtCnKGOgK3viFPbxbMTueuUUc56XsVVAvF_j0_8ZBHEH-AagSi3vyrPrFcDdt-iKDM5oCtPePcSMJePjbk82bQp8DuVV-mxvjsR2CEAsDM5yBuTUxQJyzyQFODiZJ-X0VAM4CXw0dR_JiUpzzhv-zP5GtB3snGKqKWb_saA3na76naJkOTGUFlPoNeCxeqEkDcoGHLKMotUH8Ftaxv0UBjZSeypTewmFiFnqH_wyPoIiY_dLemtQN8VXFcVPlIVEgZ0pG0Vx2B5Wn5viY_mC0)


* 新建状态机

```go
import "github.com/threeq/gofsm"

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

```

* 增加全局状态转换事件监听

这里也可以为每个事件转换单独设置事件监听

```go
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

orderStateMachine.Processor(&OrderEventProcessor{})

```

* 调用状态机，执行业务操作

```go
order := context.WithValue(context.TODO(), "data", "order object data")
state, err := orderStateMachine.Trigger(order, Start, CreateEvent)
println(fmt.Sprintf("====: %v : %v", state, err))

state, err = orderStateMachine.Trigger(order, Start, PayEvent)
println(fmt.Sprintf("====: %v : %v", state, err))

state, err = orderStateMachine.Trigger(order, Paying, PayFailureEvent)
println(fmt.Sprintf("====: %v : %v", state, err))
```

output:
```text
OnExit: [order object data] Exit [Start] on event [Create]
doAction: [order object data] --Create--> [WaitPay]
OnEnter: [order object data] Enter [WaitPay]
====: WaitPay : <nil>
====:  : 没有定义状态转换事件 [Start --Pay--> ???]
OnExit: [order object data] Exit [Paying] on event [PayFailure]
doAction: [order object data] --PayFailure--> [PayFailure]
OnEnter: [order object data] Enter [PayFailure]
====: PayFailure : <nil>
```

完整代码查看 https://github.com/threeq/gofsm/blob/master/fsm_test.go 中 `TestStateMachine_Example_Order`
