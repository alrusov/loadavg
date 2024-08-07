package loadavg

import (
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
	sync.Mutex
	firstCycle bool
	size       int
	idx        int
	lastTime   int64
	cache      []float64
	count      []int
}

//----------------------------------------------------------------------------------------------------------------------------//

// Init --
func Init(duration time.Duration) (me *LoadAvg) {
	period := int(duration / time.Second)

	if period <= 0 {
		period = 60
	} else if period > MaxPeriod {
		period = MaxPeriod
	}

	period++ // +1 - не учитываем текущую секунду

	me = &LoadAvg{
		firstCycle: true,
		size:       period,
		idx:        0,
		lastTime:   time.Now().UTC().Unix(),
		cache:      make([]float64, period),
		count:      make([]int, period),
	}

	return
}

//----------------------------------------------------------------------------------------------------------------------------//

// Add --
func (me *LoadAvg) Add(v float64) {
	me.Lock()
	defer me.Unlock()

	me.setTime()

	me.cache[me.idx] += v
	me.count[me.idx]++
}

//----------------------------------------------------------------------------------------------------------------------------//

// Value --
func (me *LoadAvg) Value() float64 {
	return me.value(false)
}

// AbsValue --
func (me *LoadAvg) AbsValue() float64 {
	return me.value(true)
}

// value --
func (me *LoadAvg) value(isAbs bool) float64 {
	me.Lock()
	defer me.Unlock()

	me.setTime()

	if me.firstCycle && me.idx == 0 {
		return 0.
	}

	v := float64(0)
	n := 0

	for i := 0; i < me.size; i++ {
		if i == me.idx {
			if me.firstCycle {
				break
			}
			continue
		}

		if isAbs {
			n += me.count[i]
		} else {
			n++
		}

		v += me.cache[i]
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
			me.count = make([]int, me.size)
			me.firstCycle = false
		} else {
			for diff > 0 {
				diff--
				me.idx++
				if me.idx == me.size {
					me.idx = 0
					me.firstCycle = false
				}
				me.cache[me.idx] = 0
				me.count[me.idx] = 0
			}
		}
	}
}

//----------------------------------------------------------------------------------------------------------------------------//
