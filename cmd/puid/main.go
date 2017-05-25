package main

import (
	"flag"
	"fmt"

	"github.com/thechriswalker/puid"
)

func example() {
	// 'c' is the default and makes this lib behave like cuid
	fmt.Println("puid.New()\n\t", puid.New())

	// specify a single byte (range is 'a-z')
	fmt.Println("puid.WithPrefixByte('z').New()\n\t", puid.WithPrefixByte('z').New())

	// like the cuid lib
	fmt.Println("puid.Cuid().New()\n\t", puid.Cuid().New())

	// prefixes can be longer and arbitrary
	// note that this really breaks the spirit of the initial cuid design
	// also you now now longer have a fixed length output, given that the prefix
	// can be fixed, this is up to you.
	fmt.Println(`puid.WithPrefix("str:").New()`+"\n\t", puid.WithPrefix("str:").New())
	fmt.Println(`puid.WithPrefixBytes([]byte("foo:")).New()`+"\n\t", puid.WithPrefixBytes([]byte("foo:")).New())
}

var (
	showExample = flag.Bool("example", false, "show example library use and output")
	count       = flag.Int("count", 1, "how many of each puid to generate")
)

func main() {
	flag.Parse()
	if *showExample {
		example()
	} else {
		prefixes := flag.Args()
		if len(prefixes) > 0 {
			for _, p := range prefixes {
				run(puid.WithPrefix(p), *count)
			}
		} else {
			run(puid.Default(), *count)
		}
	}
}

func run(g *puid.Generator, c int) {
	for i := 0; i < c; i++ {
		fmt.Println(g.New())
	}
}
