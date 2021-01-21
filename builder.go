package gosm

var stateMachines = make(map[string]*StateMachine)

type Builder struct {
	transitions []*transCfg
}

func (o *Builder) Transition() *transition {
	return &transition{&transCfg{
		builder: o,
	}}
}

func (o *Builder) Build(machineID string, options ...Option) *StateMachine {
	sm := NewMachine(options...)
	sm.Name = machineID
	for _, cfg := range o.transitions {
		for _, s1 := range cfg.from {
			sm.Trans(s1.Exit(cfg.event, cfg.condDesc, cfg.condition),
				cfg.to.Entry(cfg.actionDesc, cfg.action))
		}
	}

	stateMachines[machineID] = sm

	o.transitions = []*transCfg{}
	return sm
}

func (o *Builder) DSL(dsl string) {
	//TODO DSL 支持
}

func (o *Builder) addTransition(trans *transCfg) {
	o.transitions = append(o.transitions, trans)
}

func NewBuilder() *Builder {
	return new(Builder)
}

type transCfg struct {
	from       []*state
	to         *state
	event      Event
	condition  Condition
	action     Action
	builder    *Builder
	actionDesc string
	condDesc   string
}

type transition struct {
	*transCfg
}

type bState struct {
	*transCfg
}

type stateEvent struct {
	*transCfg
}

type stateEventWhen struct {
	*transCfg
}

type stateEventWhenAction struct {
	*transCfg
}

func (o *stateEventWhenAction) Action(desc string, action Action) *Builder {
	o.action = action
	o.actionDesc = desc
	o.builder.addTransition(o.transCfg)
	return o.builder
}

func (o *stateEventWhen) When(desc string, cond Condition) *stateEventWhenAction {
	o.condition = cond
	o.condDesc = desc
	return &stateEventWhenAction{
		transCfg: o.transCfg,
	}
}

func (o *stateEvent) On(e Event) *stateEventWhen {
	o.event = e
	return &stateEventWhen{
		transCfg: o.transCfg,
	}
}

func (o *bState) To(to *state) *stateEvent {
	o.to = to
	return &stateEvent{
		transCfg: o.transCfg,
	}
}

func (o *transition) From(from ...*state) *bState {
	o.from = from
	return &bState{
		transCfg: o.transCfg,
	}
}
