package loadavg

import (
	"errors"
	"sync"
	"time"
)

//----------------------------------------------------------------------------------------------------------------------------//

const (
	// MaxPeriod --
	MaxPeriod = 60 * 60 * 24
)

// LoadAvg --
type LoadAvg struct {
	mutex    *sync.Mutex
	cycled   bool
	size     int
	idx      int
	lastTime int64
	cache    []float64
}

//----------------------------------------------------------------------------------------------------------------------------//

// Init --
func Init(period int) (*LoadAvg, error) {
	if period <= 0 || period > MaxPeriod {
		return nil, errors.New("Bad period")
	}

	period++ // +1 - не учитываем текущую секунду

	me := &LoadAvg{
		mutex:    new(sync.Mutex),
		cycled:   false,
		size:     period,
		idx:      0,
		lastTime: time.Now().UTC().Unix(),
		cache:    make([]float64, period),
	}
	return me, nil
}

//----------------------------------------------------------------------------------------------------------------------------//

// Add --
func (me *LoadAvg) Add(v float64) {
	me.mutex.Lock()
	defer me.mutex.Unlock()

	me.setTime()

	me.cache[me.idx] += v
}

//----------------------------------------------------------------------------------------------------------------------------//

// Value --
func (me *LoadAvg) Value() float64 {
	me.mutex.Lock()
	defer me.mutex.Unlock()

	me.setTime()

	if !me.cycled && me.idx == 0 {
		return 0.
	}

	v := float64(0)
	n := me.idx

	i := 0
	barrier := me.idx
	for i < barrier {
		v += me.cache[i]
		i++
	}

	if me.cycled {
		n = me.size - 1

		i++
		barrier = me.size
		for i < barrier {
			v += me.cache[i]
			i++
		}
	}

	if n == 0 {
		return 0.
	}

	return v / float64(n)
}

//----------------------------------------------------------------------------------------------------------------------------//

func (me *LoadAvg) setTime() {
	t := time.Now().UTC().Unix()

	if t > me.lastTime {
		diff := t - me.lastTime
		me.lastTime = t
		if diff >= int64(me.size) {
			me.cache = make([]float64, me.size)
			me.cycled = true
		} else {
			for diff > 0 {
				diff--
				me.idx++
				if me.idx == me.size {
					me.idx = 0
					me.cycled = true
				}
				me.cache[me.idx] = 0
			}
		}
	}
}

//----------------------------------------------------------------------------------------------------------------------------//
