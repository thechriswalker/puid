package puid

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"gopkg.in/lucsky/cuid.v1"
)

//
// Tests
//
var prefix_str = []string{"a", "", "foo"}
var prefix_byte = []byte{'z', '*', ' '}
var prefix_bytes = [][]byte{{'a', 'b', 'c'}, {'f'}, {'x', ':'}}

func testPrefix(t *testing.T, g *Generator, prefix string) {
	id := g.New()
	if len(id) != len(prefix)+8+4*BLOCK {
		t.Errorf("unexpected length with prefix `%s`. id = %s", prefix, id)
	}
	if id[0:len(prefix)] != prefix {
		t.Errorf("unexpected prefix in id, expected `%s`, got `%s`", prefix, id[0:len(prefix)])
	}
}

func Test_Prefixes(t *testing.T) {
	for _, s := range prefix_str {
		testPrefix(t, WithPrefix(s), s)
	}
	for _, b := range prefix_byte {
		testPrefix(t, WithPrefixByte(b), string([]byte{b}))
	}
	for _, b := range prefix_bytes {
		testPrefix(t, WithPrefixBytes(b), string(b))
	}
}

// if we set our random source and counter and fingerprint up just so, we can predict what the
// algorithm should produce (we can make it pure/deterministic)
func Test_Deterministic(t *testing.T) {
	g := &Generator{
		prefix:      []byte{'x'},
		random:      badRandom(27), // 27 == "r" in base36
		fingerprint: []byte("ffff"),
		counter:     dumbCounter(1 + 36 + 36*36 + 36*36*36), // 1111
	}
	// replace the time with a fixed known value (actually the lowest possible puid value with current block size: Mon Jun 26 1972 00:49:24 GMT+0100 (BST))
	// which in milliseconds unixtime and base36 is "10000000" in milliseconds it is 78364164096 so we put it in as nano
	ft := getTime
	getTime = func() time.Time { return time.Unix(0, 78364164096*1e6) }

	// now grab the id
	actual := g.New()

	//put time back
	getTime = ft

	// prefix + time (2 blocks) + counter (1 block) + fingerprint (1 block) + random (2 blocks)
	expected := "x100000001111ffffrrrrrrrr"
	if actual != expected {
		t.Errorf("deterministic test failed. expected `%s`, got `%s`", expected, actual)
	}
}

func Test_PanicIfAppendBytesCalledWithNilSlice(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Error("we should have panic'd on a nil slice to AppendBytes")
		}
	}()
	// should do the same if we do it when we create a brand new generator
	AppendBytes(nil)
}

//
//  Benchmarks
//

func Benchmark_PuidString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New()
	}
}
func Benchmark_PuidBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Bytes()
	}
}
func Benchmark_PuidAppendBytes(b *testing.B) {
	buff := make([]byte, 9+4*BLOCK)
	for i := 0; i < b.N; i++ {
		buff = AppendBytes(buff[:])
	}
}

func Benchmark_PuidInCuidMode(b *testing.B) {
	c := Cuid()
	for i := 0; i < b.N; i++ {
		c.New()
	}
}

// just a comparison, we are faster, but not by much
func Benchmark_LucskyCuid(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cuid.New()
	}
}

//
// now some tests ripped directly from the cuid package (only slight modifications)
//

func Test_PUIDFormat(t *testing.T) {
	format := regexp.MustCompile(fmt.Sprintf("p[0-9a-z]{%d}", 8+4*BLOCK))
	if !format.MatchString(New()) {
		t.Error("Incorrect format")
	}
	if !format.Match(Bytes()) {
		t.Error("Incorrect format")
	}
}

func Test_CUIDFormat(t *testing.T) {
	c := Cuid().New()
	format := regexp.MustCompile(fmt.Sprintf("c[0-9a-z]{%d}", 8+4*BLOCK))
	if !format.MatchString(c) {
		t.Error("Incorrect format")
	}
}

func Test_SmokeTestForCollisions(t *testing.T) {
	ids := map[string]struct{}{}
	for i := 0; i < 1e6; i++ {
		id := New()
		if _, collision := ids[id]; collision {
			t.Errorf("Collision detected, at iteration %d", i)
		}
		ids[id] = struct{}{}
	}
}

func newCUID(chn chan error) {
	// use Default() here to force this code path for coverage
	Default().New()
	chn <- nil
}

func Test_DataRaces(t *testing.T) {
	chn := make(chan error)

	go newCUID(chn)
	go newCUID(chn)

	<-chn
	<-chn
}
