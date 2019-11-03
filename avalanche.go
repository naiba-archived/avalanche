package avalanche

import (
	"sync"
)

type runnerS struct {
	cond *sync.Cond
	wg   *sync.WaitGroup
}

var locker sync.Mutex
var keys map[interface{}]*runnerS

// Do something anti avalanche
func Do(key interface{}, something func() (interface{}, error)) (interface{}, error) {
	locker.Lock()
	var res interface{}
	var err error
	runner, has := keys[key]
	if !has {
		runner = &runnerS{
			sync.NewCond(new(sync.Mutex)),
			new(sync.WaitGroup),
		}
		keys[key] = runner
		locker.Unlock()
		defer func() {
			runner.wg.Wait()
			keys[key] = nil
			delete(keys, key)
			locker.Unlock()
		}()
		res, err = something()
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
	return res, err
}
