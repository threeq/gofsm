package gosm

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

var (
	ok1 = State("ok1").Entry(func(ctx context.Context, entity Entity, from, to IState) error {
		return nil
	})
	ok2 = State("ok1").Entry(func(ctx context.Context, entity Entity, from, to IState) error {
		return nil
	})
	err1 = State("err1").Entry(func(ctx context.Context, entity Entity, from, to IState) error {
		return errors.New("err1")
	})
	ok3 = State("ok1").Entry(func(ctx context.Context, entity Entity, from, to IState) error {
		return nil
	})
	err2 = State("err2").Entry(func(ctx context.Context, entity Entity, from, to IState) error {
		return errors.New("err2")
	})
)

func testBuilder(strategy SuccessStrategy, err error) func(t *testing.T) {
	return func(t *testing.T) {

	}
}

func TestAll(t *testing.T) {
	tests := []struct {
		name         string
		stateEntries []StateEntry
		want         error
	}{
		{"ok", []StateEntry{ok1, ok2, ok3}, nil},
		{"err", []StateEntry{ok1, ok2, err1, ok3}, errors.New("存在部分错误")},
		{"err2", []StateEntry{err1, err2}, errors.New("存在部分错误")},
	}
	for _, tt := range tests {
		c := context.Background()
		entity := NewTestEntity("all", State("start"))
		from := State("start")
		t.Run(tt.name+"/Parallel", func(t *testing.T) {
			got := Parallel(All)(c, entity, from, tt.stateEntries)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("All() got = %v, want %v", got, tt.want)
			}

		})
		t.Run(tt.name+"/Serial", func(t *testing.T) {
			got := Serial(All)(c, entity, from, tt.stateEntries)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("All() got = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestAllFast(t *testing.T) {
	tests := []struct {
		name         string
		stateEntries []StateEntry
		want         error
	}{
		{"ok", []StateEntry{ok1, ok2, ok3}, nil},
		{"err", []StateEntry{ok1, ok2, err1, ok3}, errors.New("err1")},
		{"err2", []StateEntry{err2}, errors.New("err2")},
	}
	for _, tt := range tests {
		c := context.Background()
		entity := NewTestEntity("all", State("start"))
		from := State("start")
		t.Run(tt.name+"/Parallel", func(t *testing.T) {
			got := Parallel(AllFast)(c, entity, from, tt.stateEntries)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("All() got = %v, want %v", got, tt.want)
			}

		})
		t.Run(tt.name+"/Serial", func(t *testing.T) {
			got := Serial(AllFast)(c, entity, from, tt.stateEntries)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("All() got = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestAlways(t *testing.T) {
	tests := []struct {
		name         string
		stateEntries []StateEntry
		want         error
	}{
		{"ok", []StateEntry{ok1, ok2, ok3}, nil},
		{"err", []StateEntry{ok1, ok2, err1, ok3}, nil},
		{"err2", []StateEntry{err1, err2}, nil},
	}
	for _, tt := range tests {
		c := context.Background()
		entity := NewTestEntity("all", State("start"))
		from := State("start")
		t.Run(tt.name+"/Parallel", func(t *testing.T) {
			got := Parallel(Always)(c, entity, from, tt.stateEntries)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("All() got = %v, want %v", got, tt.want)
			}

		})
		t.Run(tt.name+"/Serial", func(t *testing.T) {
			got := Serial(Always)(c, entity, from, tt.stateEntries)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("All() got = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestOne(t *testing.T) {
	tests := []struct {
		name         string
		stateEntries []StateEntry
		want         error
	}{
		{"ok", []StateEntry{ok1, ok2, ok3}, nil},
		{"err", []StateEntry{ok1, ok2, err1, ok3}, nil},
		{"err2", []StateEntry{err1, err2}, errors.New("全部错误")},
	}
	for _, tt := range tests {
		c := context.Background()
		entity := NewTestEntity("all", State("start"))
		from := State("start")
		t.Run(tt.name+"/Parallel", func(t *testing.T) {
			got := Parallel(One)(c, entity, from, tt.stateEntries)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("All() got = %v, want %v", got, tt.want)
			}

		})
		t.Run(tt.name+"/Serial", func(t *testing.T) {
			got := Serial(One)(c, entity, from, tt.stateEntries)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("All() got = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestOneFast(t *testing.T) {
	tests := []struct {
		name         string
		stateEntries []StateEntry
		want         error
	}{
		{"ok", []StateEntry{ok1, ok2, ok3}, nil},
		{"err", []StateEntry{ok1, ok2, err1, ok3}, nil},
		{"err2", []StateEntry{err1, err2}, errors.New("全部错误")},
	}
	for _, tt := range tests {
		c := context.Background()
		entity := NewTestEntity("all", State("start"))
		from := State("start")
		t.Run(tt.name+"/Parallel", func(t *testing.T) {
			got := Parallel(OneFast)(c, entity, from, tt.stateEntries)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("All() got = %v, want %v", got, tt.want)
			}

		})
		t.Run(tt.name+"/Serial", func(t *testing.T) {
			got := Serial(OneFast)(c, entity, from, tt.stateEntries)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("All() got = %v, want %v", got, tt.want)
			}

		})
	}
}
