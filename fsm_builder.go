package gosm

type Builder struct {
	transitions []*transCfg
}

func (owner *Builder) Transition() *transition {
	return &transition{&transCfg{
		builder: owner,
	}}
}

func (owner *Builder) Build(machineID string, options ...Option) *StateMachine {
	sm := NewStateMachine(options...)
	sm.Name = machineID
	for _, cfg := range owner.transitions {
		for _, s1 := range cfg.from {
			sm.Transition(&Transition{
				From:     s1,
				To:       cfg.to,
				Event:    cfg.event,
				CondDesc: cfg.desc,
				Cond:     cfg.condition,
				Action:   cfg.action,
			})
		}
	}

	stateMachines[machineID] = sm

	owner.transitions = []*transCfg{}
	return sm
}

func (owner *Builder) addTransition(trans *transCfg) {
	owner.transitions = append(owner.transitions, trans)
}

func NewBuilder() *Builder {
	return new(Builder)
}

type transCfg struct {
	from      []State
	to        State
	event     Event
	condition Condition
	desc      string
	action    Action
	builder   *Builder
}

type transition struct {
	*transCfg
}

type state struct {
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

func (owner *stateEventWhenAction) Action(action Action) *Builder {
	owner.action = action
	owner.builder.addTransition(owner.transCfg)
	return owner.builder
}

func (owner *stateEventWhen) When(desc string, cond Condition) *stateEventWhenAction {
	owner.condition = cond
	owner.desc = desc
	return &stateEventWhenAction{
		transCfg: owner.transCfg,
	}
}

func (owner *stateEvent) On(e Event) *stateEventWhen {
	owner.event = e
	return &stateEventWhen{
		transCfg: owner.transCfg,
	}
}

func (owner *state) To(to State) *stateEvent {
	owner.to = to
	return &stateEvent{
		transCfg: owner.transCfg,
	}
}

func (owner *transition) From(from ...State) *state {
	owner.from = from
	return &state{
		transCfg: owner.transCfg,
	}
}

