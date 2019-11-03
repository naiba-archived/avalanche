package avalanche

import (
	"sync"
)

type runnerS struct {
	cond *sync.Cond
	wg   *sync.WaitGroup
	res  interface{}
	err  error
}

var locker sync.Mutex
var keys map[interface{}]*runnerS

func init() {
	keys = make(map[interface{}]*runnerS)
}

// Do something anti avalanche
func Do(uniqueKey interface{}, something func() (interface{}, error)) (interface{}, error) {
	locker.Lock()
	runner, has := keys[uniqueKey]
	if !has {
		runner = &runnerS{
			sync.NewCond(new(sync.Mutex)),
			new(sync.WaitGroup),
			nil,
			nil,
		}
		keys[uniqueKey] = runner
		locker.Unlock()
		defer func() {
			runner.wg.Wait()
			keys[uniqueKey] = nil
			delete(keys, uniqueKey)
			locker.Unlock()
		}()
		runner.res, runner.err = something()
		locker.Lock()
		runner.cond.L.Lock()
		runner.cond.Broadcast()
		runner.cond.L.Unlock()
	} else {
		runner.cond.L.Lock()
		runner.wg.Add(1)
		locker.Unlock()
		defer func() {
			runner.wg.Done()
			runner.cond.L.Unlock()
		}()
		runner.cond.Wait()
	}
	return runner.res, runner.err
}
