package gosm

type TestEntity struct {
	id string
	s  IState
}

func (t *TestEntity) ID() string {
	return t.id
}

func (t *TestEntity) State() IState {
	return t.s
}

func NewTestEntity(id string, s  IState) *TestEntity {
	return &TestEntity{id, s}
}

