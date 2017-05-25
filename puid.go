package puid

import (
	"strconv"
	"time"
)

const (
	BLOCK        = 4
	BASE         = 36
	MAX_HALF_INT = 1296    // 36^2 (half block size)
	MAX_COUNTER  = 1679616 //36^4 (full block size)
)

// this is just so we can be deterministic during testing
// by replacing this function
var getTime = func() time.Time {
	return time.Now()
}

func hammertime() int64 {
	return getTime().UnixNano() / int64(time.Millisecond)
}

// An puid generator
type Generator struct {
	fingerprint []byte
	random      Random
	counter     Counter
	prefix      []byte
}

// These are the options you can customize should you want
type Options struct {
	Fingerprint []byte  // this is the fingerprint for this host
	Random      Random  // this is the source of randomness
	Counter     Counter // this is the increasing counter
	Prefix      []byte  // this is the prefix ("c" in `cuid`)
}

// Spit out a new puid from the generator, raw bytes
func (g *Generator) Bytes() []byte {
	// the buffer is going to be about len(prefix) + 6*BLOCK long
	// that depends on how big a timestamp can get.
	// Timestamps will be 9 digits of base36 after: Fri Apr 22 5188 12:04:28 GMT+0100 (BST)
	// so we can probably ignore that.
	// also they were only 7 digits or les before: Mon Jun 26 1972 00:49:24 GMT+0100 (BST)
	// so we can pretty much guarrantee that the length of an ID
	// is prefix + 8 + 4*BLOCK
	buff := make([]byte, 0, len(g.prefix)+8+4*BLOCK)
	// set the prefix
	buff = append(buff, g.prefix...)
	// timestamp is not padded an 8 digits in all likelyhood (see previous comment)
	buff = strconv.AppendInt(buff, hammertime(), BASE)
	// now the counter
	buff = appendPaddedInt(buff, g.counter.Next(), BLOCK)
	// then the fingerprint (we clamped it to BLOCK bytes)
	buff = append(buff, g.fingerprint...)
	// now the random data
	buff = appendRandomBase36(buff, g.random, BLOCK*2)
	return buff
}

// Returns the raw byte slice of an puid
func Bytes() []byte {
	return defaultGenerator.Bytes()
}

// Generate an puid as a string
func (g *Generator) New() string {
	return string(g.Bytes())
}

// Returns an puid from the default generator
func New() string {
	return defaultGenerator.New()
}

// Create a clone of this generator but with the given Counter
func (g *Generator) WithCounter(c Counter) *Generator {
	if c == nil {
		panic("WithCounter called with nil Counter")
	}
	return &Generator{
		fingerprint: g.fingerprint,
		counter:     c,
		random:      g.random,
		prefix:      g.prefix,
	}
}

// create an id generator from the default one, but with the given Counter
func WithCounter(c Counter) *Generator {
	return defaultGenerator.WithCounter(c)
}

// Return a new generator like this one, but with a different source
// of randomness
func (g *Generator) WithRandom(r Random) *Generator {
	if r == nil {
		panic("WithRandom called with nil Random")
	}
	return &Generator{
		fingerprint: g.fingerprint,
		counter:     g.counter,
		random:      r,
		prefix:      g.prefix,
	}
}

// Returns a clone of the default generator with the given Random-ness source
func WithRandom(r Random) *Generator {
	return defaultGenerator.WithRandom(r)
}

// Return a new generator like this one, but with a different prefix
// remember that cuid's a supposed to be portable/url safe/start with 'a-z'
func (g *Generator) WithPrefixBytes(prefix []byte) *Generator {
	return &Generator{
		fingerprint: g.fingerprint,
		counter:     g.counter,
		random:      g.random,
		prefix:      prefix,
	}
}

// Returns the default generator but with the given []byte prefix
func WithPrefixBytes(b []byte) *Generator {
	return defaultGenerator.WithPrefixBytes(b)
}

// Return a new generator like this one, but with a different prefix given as a string
func (g *Generator) WithPrefix(s string) *Generator {
	return g.WithPrefixBytes([]byte(s))
}

// Returns the default generator but with the given string prefix
func WithPrefix(s string) *Generator {
	return defaultGenerator.WithPrefix(s)
}

// Return a new generator like this one with the single byte prefix
// the whole point of cuid is that this byte is 'a-z'
func (g *Generator) WithPrefixByte(char byte) *Generator {
	return g.WithPrefixBytes([]byte{char})
}

// Returns the default generator but with the given byte prefix
func WithPrefixByte(b byte) *Generator {
	return defaultGenerator.WithPrefixByte(b)
}

// Create a generator with the given fingerprint
// panics if the byte slice contains non-base36 characters
func (g *Generator) WithFingerprintBytes(b []byte) *Generator {
	if b == nil {
		panic("*(puid.Generator).WithFingerprintBytes called with nil byte slice")
	}
	b = massageFingerprint(b)
	return &Generator{
		fingerprint: b,
		counter:     g.counter,
		random:      g.random,
		prefix:      g.prefix,
	}
}

// Returns the default generator but with the given fingerprint []byte
// panics if the byte slice contains non-base36 characters
func WithFingerprintBytes(b []byte) *Generator {
	return defaultGenerator.WithFingerprintBytes(b)
}

// Create a generator with a fingerprint created from these values
func (g *Generator) WithFingerprint(str string, num int64) *Generator {
	fp := CreateFingerprint(str, num)
	return g.WithFingerprintBytes(fp)
}

// Returns the default generator but with a fingerprint created from the given values
func WithFingerprint(str string, num int64) *Generator {
	return defaultGenerator.WithFingerprint(str, num)
}

// these are the default options
var (
	defaultPrefix      = []byte{'p'}
	defaultFingerprint []byte
	defaultGenerator   *Generator
)

// we do all initializations here, or non-determinism means we can't know what
// order the inits from other files are

func init() {
	defaultFingerprint = getDefaultFingerprint()

	//New generator uses defaults for evrything
	defaultGenerator = NewGenerator(nil)
}

// Returns a regular cuid generator
func Cuid() *Generator {
	return NewGenerator(&Options{Prefix: []byte{'c'}})
}

// Access the default generator (in case you want to pass it around)
func Default() *Generator {
	return defaultGenerator
}

// Create a new puid generator
func NewGenerator(o *Options) *Generator {
	if o == nil {
		return &Generator{
			fingerprint: clone(defaultFingerprint),
			random:      getDefaultRandom(),
			counter:     getDefaultCounter(),
			prefix:      clone(defaultPrefix),
		}
	}
	g := &Generator{
		fingerprint: o.Fingerprint,
		random:      o.Random,
		counter:     o.Counter,
		prefix:      o.Prefix,
	}
	g.fingerprint = massageFingerprint(g.fingerprint)

	if g.random == nil {
		g.random = getDefaultRandom() // note we call this again to ensure it is a *new* random source
	}
	if g.counter == nil {
		g.counter = getDefaultCounter() // note we call this again to get a NEW counter
	}
	if g.prefix == nil {
		g.prefix = clone(defaultPrefix)
	}
	return g
}

func massageFingerprint(fp []byte) []byte {
	if fp == nil || len(fp) == 0 {
		fp = clone(defaultPrefix)
	}
	for len(fp) < BLOCK {
		//right pad is easier...
		fp = append(fp, '0')
	}
	if len(fp) > BLOCK {
		fp = fp[0:BLOCK]
	}
	if !isAllBase36(fp) {
		panic("supplied fingerprint is not base36, did you use `puid.CreateFingerprint(str, int)`?")
	}
	return fp
}

func clone(b []byte) []byte {
	a := make([]byte, len(b))
	copy(a[:], b)
	return a
}
