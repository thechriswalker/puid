package puid

import (
	"math/rand"
	"sync"
	"time"
)

// the default random is the one here
func getDefaultRandom() Random {
	return NewMathRandom(nil)
}

// The randomness source should fill the buffer
// with random data. I initially made an implementation that used
// `crypto/rand` but it was waaaay slow and unnecessary
type Random interface {
	Read([]byte) (int, error)
}

type mathRandom struct {
	r   *rand.Rand
	mtx sync.Mutex //required for concurrent access to the source
}

// A new source of psuedo-randomness for our Generator
func NewMathRandom(source rand.Source) Random {
	if source == nil {
		source = rand.NewSource(time.Now().UnixNano())
	}
	return &mathRandom{r: rand.New(source)}
}

// implements the Random interface
func (m *mathRandom) Read(b []byte) (n int, e error) {
	m.mtx.Lock()
	n, e = m.r.Read(b)
	m.mtx.Unlock()
	return
}

// append and return
func appendRandomBase36(b []byte, r Random, count int) []byte {
	t := make([]byte, count)
	r.Read(t)
	base36convert(t)
	return append(b, t...)
}
