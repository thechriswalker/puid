[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](http://godoc.org/github.com/thechriswalker/puid) [![Build Status](https://travis-ci.org/thechriswalker/puid.svg?branch=master)](https://travis-ci.org/thechriswalker/puid) 

# `puid`, like [`cuid`](https://github.com/ericelliott/cuid) (Collision-Resistant IDs)

**Please refer to the [main project site](http://usecuid.org) for the full rationale behind CUIDs.**

Almost all of the rationale stands except than with an puid, you can vary the initial character (or prefix) -- it's not just a fixed `c`.

My use case is that I want to use that character as a type hint. Purely
a visual thing. 

NB You can easily break the URL-Safety that CUIDs guarrantee by using stupid prefixes. I define this as user error, the puid library makes no attempt to stop you from doing this.

Otherwise the library and implementation is compliant and marginally faster than the other Go implementation. Thats probably because I don't use exactly the same alogorithms for generating the data. But given that most of it is "random"

<details><summary>subjective go benchmark of this lib vs. [lucsky/cuid](https://github.com/lucsky/cuid)</summary>

#### Run on my laptop, pinch of salt necessary

```
$ go test -run=XXX -bench=.
Benchmark_PuidString         3000000           411 ns/op
Benchmark_PuidBytes          5000000           362 ns/op
Benchmark_PuidAppendBytes    5000000           344 ns/op
Benchmark_PuidInCuidMode     3000000           402 ns/op
Benchmark_LucskyCuid         3000000           564 ns/op
PASS
ok      github.com/thechriswalker/puid  9.842s
```

Event more reason to question this is that technically `Benchmark_PuidString` and `Benchmark_PuidInCuidMode` are doing exactly the same work...
</details>

## install

```
$ go get github.com/thechriswalker/puid/...
```



## use as a tool

If your `$GOPATH`/`$GOBIN` is all setup, then you should be able to run

```
$ puid
pj34aln7t0000cars6wqeent6
$ puid d
dj34alh5d0000carshj2ctw8e
$ puid foo: bar:
foo:j34al93s0000carskgw1gt3q
bar:j34al93s0001carskzagqwy2
$ puid -count=10
pj34am4tg0000carsa3zp2gqh
pj34am4th0001cars809dw0k0
pj34am4th0002carsgh2xlo91
pj34am4th0003carsnrt75ob7
pj34am4th0004cars9v53axm5
pj34am4th0005carsnmzjbaoj
pj34am4ti0006carsd2wxrfks
pj34am4ti0007cars32qh9yq9
pj34am4ti0008carsezmwei9e
pj34am4ti0009carsm5pziq9m
``` 

## use as a lib

docs [on godoc.org](https://godoc.org/github.com/thechriswalker/puid)

or see [`cmd/puid/main.go`](cmd/puid/main.go)

or run `puid -example`

## customization the generator

More customization can be made to the id generator, you can supply a custom fingerprint, random data source, or counter as well as the prefix. see the docs for info

## acknowledgements

Firstly, [ericelliott](https://github.com/ericelliott) for CUIDs which are awesome.
Secondly, [lucsky](https://github.com/lucsky) for his [implementation of CUID in golang](https://github.com/lucsky/cuid) which inspired this.
