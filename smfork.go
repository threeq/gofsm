package gosm

import (
	"context"
	"errors"
	"log"
	"sync"
	"sync/atomic"
)

func Fork(executor Executor, stateEntries ...StateEntry) StateEntry {
	return &forkStateEntry{
		stateEntries: stateEntries,
		executor:     executor,
	}
}

type SuccessStrategy = func(ctx context.Context, entity Entity, from IState, stateEntries []StateEntry) (func() error, func(err error) bool)
type Executor = func(ctx context.Context, entity Entity, from IState, stateEntries []StateEntry) error

type forkStateEntry struct {
	state        IState
	action       Action
	desc         string
	stateEntries []StateEntry
	executor     Executor
}

func (f *forkStateEntry) State() IState {
	if f.state == nil {
		f.state = State("Fork")
	}
	return f.state
}

func (f *forkStateEntry) Action(ctx context.Context, entity Entity, from, to IState) error {
	return f.executor(ctx, entity, from, f.stateEntries)
}

func (f *forkStateEntry) Desc() string {
	// TODO desc
	return f.desc
}

func Serial(successStrategy SuccessStrategy) Executor {
	return func(ctx context.Context, entity Entity, from IState, stateEntries []StateEntry) error {
		wait, checker := successStrategy(ctx, entity, from, stateEntries)
		for _, entry := range stateEntries {
			err := entry.Action(ctx, entity, from, entry.State())
			stop := checker(err)
			if stop {
				break
			}
		}
		return wait()
	}
}

func Parallel(successStrategy SuccessStrategy) Executor {
	return func(ctx context.Context, entity Entity, from IState, stateEntries []StateEntry) error {
		wait, checker := successStrategy(ctx, entity, from, stateEntries)
		for _, entry := range stateEntries {
			go func(entry StateEntry) {
				//FIXME 思考是否需要前置检测
				err := entry.Action(ctx, entity, from, entry.State())
				_ = checker(err)
			}(entry)
		}
		return wait()
	}
}

//OneFast 快速成功。成功以后的 entry 不会触发
func OneFast(ctx context.Context, entity Entity, from IState, stateEntries []StateEntry) (func() error, func(err error) bool) {

	wait := &sync.WaitGroup{}
	wait.Add(1)
	once := &sync.Once{}

	var total = int32(len(stateEntries))
	var errCount = int32(0)

	var waiter = func() error {
		wait.Wait()
		if errCount == total {
			return errors.New("全部错误")
		}
		return nil
	}

	var checker = func(err error) bool {
		if err == nil {
			once.Do(func() {
				wait.Done()
			})
			return true
		}

		//FIXME 错误处理
		log.Printf("%v", err)
		atomic.AddInt32(&errCount, 1)
		if atomic.LoadInt32(&errCount) == total {
			once.Do(func() {
				wait.Done()
			})
		}

		return false
	}

	return waiter, checker
}

//One 所有 entry 都会触发
func One(ctx context.Context, entity Entity, from IState, stateEntries []StateEntry) (func() error, func(err error) bool) {

	wait := &sync.WaitGroup{}
	wait.Add(len(stateEntries))
	var success = false
	var waiter = func() error {
		wait.Wait()
		if success {
			return nil
		}
		return errors.New("全部错误")
	}

	var checker = func(err error) bool {
		if err == nil {
			success = true
		}
		//FIXME 错误处理
		if err != nil {
			log.Printf("%v", err)
		}
		wait.Done()
		return false
	}

	return waiter, checker
}

//AllFast 快速失败。失败以后的 entry 不会触发
func AllFast(ctx context.Context, entity Entity, from IState, stateEntries []StateEntry) (func() error, func(err error) bool) {
	wait := &sync.WaitGroup{}
	wait.Add(1)
	once := &sync.Once{}

	var total = int32(len(stateEntries))
	var okCount = int32(0)
	var hasErr error

	var waiter = func() error {
		wait.Wait()
		return hasErr
	}

	var checker = func(err error) bool {
		if err != nil {
			hasErr = err
			log.Printf("%v", err)
			once.Do(func() {
				wait.Done()
			})
			return true
		}
		atomic.AddInt32(&okCount, 1)
		if atomic.LoadInt32(&okCount) == total {
			once.Do(func() {
				wait.Done()
			})
		}
		return false
	}

	return waiter, checker
}

//All 所有 entry 都会触发
func All(ctx context.Context, entity Entity, from IState, stateEntries []StateEntry) (func() error, func(err error) bool) {
	wait := &sync.WaitGroup{}
	wait.Add(len(stateEntries))

	hasErr := false

	var waiter = func() error {
		wait.Wait()
		if hasErr {
			return errors.New("存在部分错误")
		}
		return nil
	}

	var checker = func(err error) bool {
		//FIXME 错误处理
		if err != nil {
			hasErr = true
			log.Printf("%v", err)
		}
		wait.Done()
		return false
	}

	return waiter, checker
}

//Always 始终正常，保证触发所有 entry
func Always(ctx context.Context, entity Entity, from IState, stateEntries []StateEntry) (func() error, func(err error) bool) {

	var waiter = func() error {
		return nil
	}

	var checker = func(err error) bool {
		//FIXME 错误处理
		if err != nil {
			log.Printf("%v", err)
		}
		return false
	}

	return waiter, checker
}
