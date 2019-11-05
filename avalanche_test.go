package avalanche

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

var i int64

func TestDo(t *testing.T) {
	var wg sync.WaitGroup
	for j := 0; j < 200; j++ {
		wg.Add(1)
		go t.Run(fmt.Sprintf("Runner %d", j), func(t1 *testing.T) {
			defer wg.Done()
			t.Log("Runner", j, "want to do something")
			res, err := Do("test", something(t))
			if err != nil || res.(int64) > 1 {
				t1.FailNow()
				return
			}
			t.Log("Runner", j, "getted res", res, err)
		})
	}
	wg.Wait()
}

func something(t *testing.T) func() (interface{}, error) {
	return func() (interface{}, error) {
		i++
		t.Log("Do something", i)
		time.Sleep(time.Second * 3)
		t.Log("Return something", i)
		return i, nil
	}
}
