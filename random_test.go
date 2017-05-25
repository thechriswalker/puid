package puid

import (
	"testing"
)

type badRandom byte

func (br badRandom) Read(b []byte) (n int, err error) {
	for i := range b {
		b[i] = byte(br)
		n++
	}
	return
}

var ones = badRandom('\x01')

func Test_CustomRandom(t *testing.T) {
	g := WithRandom(ones)
	// the final two blocks should always be ones
	r := "11111111"
	confirmDumbRandom(t, g, r)

	// should be the same if create a new generator from scratch
	g = NewGenerator(&Options{Random: ones})
	confirmDumbRandom(t, g, r)
}

func confirmDumbRandom(t *testing.T, g *Generator, r string) {
	id := g.New()
	if id[9+BLOCK*2:] != r {
		t.Errorf("unexpected random data section in id: %s", id)
	}
	//and again
	id = g.New()
	if id[9+BLOCK*2:] != r {
		t.Errorf("unexpected random data section in id: %s", id)
	}
}

func Test_NilRandomCausesPanic(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Error("we should have panic'd on a nil Random")
		}
	}()
	WithRandom(nil)
}
