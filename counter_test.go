package puid

import (
	"bytes"
	"strconv"
	"testing"
)

func Test_CounterRollover(t *testing.T) {
	// a counter starting in the middle
	ctr := &counterMutex{value: MAX_COUNTER / 2}

	// lets do 3 full loops
	n := MAX_COUNTER * 3

	var c int64
	for i := 0; i < n; i++ {
		c = ctr.Next()
		if c >= MAX_COUNTER {
			t.Fatalf("counter exceeded MAX (iteration: %d, counter: %d)", i, c)
		}
	}
}

type dumbCounter int64

func (d dumbCounter) Next() int64 {
	return int64(d)
}

func Test_CustomCounter(t *testing.T) {
	c := dumbCounter(1337) // => "0115" in padded base36
	g := WithCounter(c)
	id := g.New()
	// the slice of this at should be the counter value
	if id[9:9+BLOCK] != "0115" {
		t.Fatalf("unexpected custom counter value in puid: %s", id)
	}
	// and the next one should be the same
	id = g.New()
	if id[9:9+BLOCK] != "0115" {
		t.Fatalf("unexpected custom counter value in puid: %s", id)
	}
}

func Test_NilCounterCausesPanic(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Error("we should have panic'd on a nil Counter")
		}
	}()
	WithCounter(nil)
}

func Test_Strlen(t *testing.T) {
	tests := []struct {
		n int64
		l int
	}{
		{0, 1},
		{1, 1},
		{35, 1},
		{36, 2},
		{36 * 36, 3},
		{36 * 36 * 36, 4},
		{36*36*36 + 100, 4},
		{36*36*36 - 1, 3},
		{36 * 36 * 36 * 36 * 36 * 36, 6},
	}

	for _, test := range tests {
		actual := strlen(test.n)
		if actual != test.l {
			t.Errorf("unexpected strlen for %s (%d), expected: %d, got %d", strconv.FormatInt(test.n, BASE), test.l, actual)
		}
	}
}

func Test_AppendPaddedInt(t *testing.T) {
	tests := []struct {
		n int64
		e string
	}{
		{0, "0000"},
		{1, "0001"},
		{35, "000z"},
		{36, "0010"},
		{36 * 36, "0100"},
		{36 * 36 * 36, "1000"},
		{36*36*36 + 100, "102s"},
		{36*36*36 - 1, "0zzz"},
		{36 * 36 * 36 * 36 * 36 * 36, "1000000"},
	}

	for _, tt := range tests {
		b := make([]byte, 0, 4)
		b = appendPaddedInt(b, tt.n, 4)
		if !bytes.Equal(b, []byte(tt.e)) {
			t.Errorf("unexpected string padding, \n\tgot `%x`\n\texp `%x`", b, tt.e)
		}
	}
}
