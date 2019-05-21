# TDD and Benchmarking in Go

I was experimenting with [TinyGo](http://github.com/mramshaw/TinyGo) and thought I'd take a stab
at string reversal. This is pretty easy to do __bytewise__ in TinyGo but not yet possible __runewise__
(rune conversion is on the TinyGo todo list).

Following discussions about how to reverse a string, a quick search revealed the following code:

    http://github.com/golang/example/tree/master/stringutil

And a podcast I was listening to discussed how optimization goals might change depending upon use case.

Accordingly, it seemed to be time for a refresher on TDD and Benchmarking in Golang using
this code as a basis.

The stated goals for the code are:

* The basic form of a library
* Conversion between string and []rune
* Table-driven unit tests ([testing](http://golang.org/pkg/testing/))

My intent here is to examine (and possibly improve) the third goal.

## Contents

The contents are as follows:

* [Prequisites](#prerequisites)
* [stringutil](#stringutil)
    * [Current tests](#current-tests)
    * [Current code coverage](#current-code-coverage)
    * [Current benchmarks](#current-benchmarks)
* [Discussion](#discussion)
    * [Byte](#byte)
    * [Rune](#rune)
    * [Strings](#strings)
* [Reference](#reference)
* [To Do](#to-do)
* [Credits](#credits)

## Prerequisites

Golang (version __1.11__ or better) installed and configured.

## stringutil

Rather than create some new code, by starting with an existing small codebase we
can specifically restrict our efforts to testing and benchmarking. This approach
has its advantages.

#### Current tests

As things stand, only one test (TestReverse) is currently defined:

```bash
$ cd stringutil
$ go test -v
=== RUN   TestReverse
--- PASS: TestReverse (0.00s)
PASS
ok  	_/home/owner/Documents/GO/TDD_and_Benchmarking_in_Go/stringutil	0.001s
$
```

However, this is a ___table-driven test___ which covers multiple tests (including, most importantly,
the ___empty string___ - which is often a useful edge case).

#### Current code coverage

Currently, 100% of the code is covered (this is generally only possible for small libraries):

```bash
$ go test -v -cover
=== RUN   TestReverse
--- PASS: TestReverse (0.00s)
PASS
coverage: 100.0% of statements
ok  	_/home/owner/Documents/GO/TDD_and_Benchmarking_in_Go/stringutil	0.001s
$
```

#### Current benchmarks

After coding up some quick benchmarks, the results are:

```bash
$ go test -v -bench=. -cover
=== RUN   TestReverse
--- PASS: TestReverse (0.00s)
goos: linux
goarch: amd64
BenchmarkReverseBytes-4         	10000000	       135 ns/op
BenchmarkReverseEmptyString-4   	50000000	        28.6 ns/op
BenchmarkReverseRunes-4         	10000000	       135 ns/op
PASS
coverage: 100.0% of statements
ok  	_/home/owner/Documents/GO/TDD_and_Benchmarking_in_Go/stringutil	4.440s
$
```

[These results are for __my__ computer; your results will probably be different.]

## Discussion

Some points to keep in mind are as follows.

[These are from the [builtin](http://golang.org/pkg/builtin/) package documentation, which is
 where native Golang data types - which would otherwise be undocumented - are documented.]

#### Byte

From:

    http://golang.org/pkg/builtin/#byte

> byte is an alias for uint8 and is equivalent to uint8 in all ways.
> It is used, by convention, to distinguish byte values from 8-bit unsigned integer values.

#### Rune

From:

    http://golang.org/pkg/builtin/#rune

> rune is an alias for int32 and is equivalent to int32 in all ways.
> It is used, by convention, to distinguish character values from integer values.

#### Strings

As in Java, strings are immutable. This has consequences for memory allocation as well
as garbage collection (at the moment, TinyGo has minimal GC).

From:

    http://golang.org/pkg/builtin/#string

> string is the set of all strings of 8-bit bytes, conventionally but not necessarily
> representing UTF-8-encoded text. A string may be empty, but not nil.
> Values of string type are immutable.

## Reference

Back in the day, text processing in ASCII was relatively straightforward.

But ASCII only really covered English, and there are many other languages.

Most of the European languages were (somewhat) easily dealt-with, but the
CJKV (Chinese, Japanese, Korean and Vietnamese) languages presented some
very difficult edge cases. Accordingly, we have ___runes___.

The Go Blog has some interesting articles. The following may be useful:

    http://blog.golang.org/strings

And:

    http://blog.golang.org/normalization

## To Do

- [x] Implement benchmarks
- [ ] Implement better tests
- [ ] Create different implementations of the code & benchmark them
- [ ] Verify the optimization criteria for TinyGo
- [ ] Implement __runes__ in TinyGo

## Credits

Somewhat inspired by this fascinating podcast ("Compilers, LLVM, Swift, TPU, and ML Accelerators"):

    http://lexfridman.com/chris-lattner/

[Lex Fridman usually has fascinating guests & prepares really interesting questions.
 The key reminder was that it is possible to optimize for different use cases.]

Also by this interesting Go Time podcast ("It's time to talk about testing"):

    http://changelog.com/gotime/83

[The inspiration here was the discussion around whether or not 100% code
 coverage was really sufficient. TL;DR - it really depends upon the test
 quality. My personal view is that 100% code coverage is not always
 desireable or even possible. Nevertheless, more code coverage is
 almost always (but not always) better than less code coverage.]
