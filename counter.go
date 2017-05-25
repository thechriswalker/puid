package puid

import (
	"strconv"
	"sync"
)

// This interface is all that's needed to be a Counter
// for the puid's.
// You probably don't want to create one of these unless
// you are trying to co-ordinate something like using a
// redis instance for atomic incrememt counters across machines.
// If you do, do not let the value returned b < 0 or > MAX_COUNTER
type Counter interface {
	Next() int64
}

// a simple mutex protected counter
// I used a seperate interface to test this vs. a channel
// the channel, was an order of magnitude slower.
type counterMutex struct {
	value int64
	mtx   sync.Mutex
}

func (c *counterMutex) Next() (n int64) {
	c.mtx.Lock()
	n, c.value = c.value, c.value+1
	if c.value == MAX_COUNTER {
		c.value = 0
	}
	c.mtx.Unlock()
	return
}

func getDefaultCounter() *counterMutex {
	return &counterMutex{}
}

// This is a single Block size of 4 assumed
func appendPaddedInt(b []byte, v int64, size int) []byte {
	d := size - strlen(v)
	switch d {
	case 1:
		b = append(b, '0')
	case 2:
		b = append(b, '0', '0')
	case 3:
		b = append(b, '0', '0', '0')
	default:
		if d > 0 {
			for i := 0; i < d; i++ {
				b = append(b, '0')
			}
		}
	}
	// then write the converted int
	return strconv.AppendInt(b, v, BASE)

}

const (
	_1_DIGIT = BASE
	_2_DIGIT = _1_DIGIT * BASE
	_3_DIGIT = _2_DIGIT * BASE
	_4_DIGIT = _3_DIGIT * BASE
)

// string length of i in base 36
// i.e. if we write i in base 36, how many digits do we use
func strlen(i int64) (l int) {
	switch {
	case i < _1_DIGIT:
		return 1
	case i < _2_DIGIT:
		return 2
	case i < _3_DIGIT:
		return 3
	case i < _4_DIGIT:
		return 4
	default:
		//now we have to work it out... (bnut we only have to start from 4)
		l = 4
		i = i / _4_DIGIT
		for _1_DIGIT < i {
			l++
			i = i / _1_DIGIT
		}
		l++
		return
	}
}
