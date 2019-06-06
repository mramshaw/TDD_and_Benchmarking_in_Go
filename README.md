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

My intent here is to examine (and possibly improve) the third goal. In addition, I will seek to
optimize the code for different use cases (while examining these use cases).

## Contents

The contents are as follows:

* [Prequisites](#prerequisites)
* [stringutil](#stringutil)
    * [Current tests](#current-tests)
    * [Current code coverage](#current-code-coverage)
    * [Current benchmarks](#current-benchmarks)
    * [Memory allocations](#memory-allocations)
    * [Garbage Collection](#garbage-collection)
    * [Alternatives](#alternatives)
* [Discussion](#discussion)
    * [Byte](#byte)
    * [Rune](#rune)
    * [Strings](#strings)
* [Reference](#reference)
* [Versions](#versions)
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

However, this is a [table-driven test](http://github.com/mramshaw/radix-trie#table-driven-tests)
which covers multiple tests (including, most importantly, the ___empty string___ - which is
often a useful edge case).

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

#### Memory allocations

The above results are sufficient for comparing different implementations of string reversal,
however - just for fun - lets explore the memory allocations:

```bash
$ go test -v -bench=. -cover -benchmem
=== RUN   TestReverse
--- PASS: TestReverse (0.00s)
goos: linux
goarch: amd64
BenchmarkReverseBytes-4         	10000000	       135 ns/op	      16 B/op	       1 allocs/op
BenchmarkReverseEmptyString-4   	50000000	        28.5 ns/op	       3 B/op	       1 allocs/op
BenchmarkReverseRunes-4         	10000000	       134 ns/op	      16 B/op	       1 allocs/op
PASS
coverage: 100.0% of statements
ok  	_/home/owner/Documents/GO/TDD_and_Benchmarking_in_Go/stringutil	4.431s
$
```

Here we can see that the overall benchmarks have not changed significantly but now we see both
the number of memory allocations as well as the size of these allocations (in Bytes). Both of
these are average values, where the values are averaged over the time the benchmarks were run.

We can get more atomic information about our Benchmark Results if we wish:

    http://golang.org/pkg/testing/#BenchmarkResult

Lets have another look at our benchmarks to see how they differ:

```Golang
func BenchmarkReverseBytes(b *testing.B) {
	benchmarkReverse("Hello, world", b)
}

func BenchmarkReverseEmptyString(b *testing.B) {
	benchmarkReverse("", b)
}

func BenchmarkReverseRunes(b *testing.B) {
	benchmarkReverse("Hello, 世界", b)
}
```

#### Garbage Collection

When profiling code, Garbage Collection needs to be taken into account - as it can badly skew
comparative benchmarks. So lets see with what frequency GC is actually occurring:

```bash
$ GODEBUG=gctrace=1 go test -v -bench=. -cover -benchmem
gc 1 @0.021s 1%: 0.013+1.6+0.19 ms clock, 0.052+0.21/0.38/0.56+0.79 ms cpu, 4->4->0 MB, 5 MB goal, 4 P
<snip>
gc 9 @0.082s 4%: 1.5+8.1+0.032 ms clock, 6.2+0.23/2.6/0.28+0.13 ms cpu, 5->5->1 MB, 7 MB goal, 4 P
# _/home/owner/Documents/GO/TDD_and_Benchmarking_in_Go/stringutil.test
gc 1 @0.000s 17%: 0.002+2.4+0.017 ms clock, 0.009+0.094/2.3/0.12+0.070 ms cpu, 4->4->4 MB, 5 MB goal, 4 P
<snip>
gc 4 @0.059s 7%: 0.004+6.9+0.042 ms clock, 0.018+0.063/6.7/7.2+0.16 ms cpu, 26->28->26 MB, 28 MB goal, 4 P
=== RUN   TestReverse
--- PASS: TestReverse (0.00s)
gc 1 @0.000s 9%: 0.015+0.063+0.022 ms clock, 0.062+0/0.026/0.083+0.091 ms cpu, 0->0->0 MB, 4 MB goal, 4 P (forced)
gc 2 @0.000s 13%: 0.019+0.031+0.014 ms clock, 0.077+0/0.027/0.058+0.057 ms cpu, 0->0->0 MB, 4 MB goal, 4 P (forced)
goos: linux
goarch: amd64
BenchmarkReverseBytes-4         	gc 3 @0.000s 12%: 0.001+0.039+0.008 ms clock, 0.007+0/0.023/0.063+0.034 ms cpu, 0->0->0 MB, 4 MB goal, 4 P (forced)
gc 4 @0.000s 12%: 0.001+0.028+0.008 ms clock, 0.005+0/0.023/0.053+0.032 ms cpu, 0->0->0 MB, 4 MB goal, 4 P (forced)
gc 5 @0.002s 5%: 0.001+0.041+0.008 ms clock, 0.005+0/0.025/0.062+0.033 ms cpu, 0->0->0 MB, 4 MB goal, 4 P (forced)
gc 6 @0.037s 0%: 0.002+0.073+0.041 ms clock, 0.010+0.056/0.016/0.044+0.16 ms cpu, 4->4->0 MB, 5 MB goal, 4 P
gc 7 @0.072s 0%: 0.002+0.076+0.038 ms clock, 0.010+0.046/0.019/0.037+0.15 ms cpu, 4->4->0 MB, 5 MB goal, 4 P
gc 8 @0.107s 0%: 0.002+0.083+0.042 ms clock, 0.011+0.046/0.021/0.036+0.16 ms cpu, 4->4->0 MB, 5 MB goal, 4 P
gc 9 @0.139s 0%: 0.018+0.046+0.016 ms clock, 0.074+0/0.041/0.054+0.067 ms cpu, 3->3->0 MB, 4 MB goal, 4 P (forced)
<snip>
gc 48 @1.486s 0%: 0.002+0.063+0.022 ms clock, 0.009+0.044/0.023/0.032+0.089 ms cpu, 4->4->0 MB, 5 MB goal, 4 P
10000000	       135 ns/op	      16 B/op	       1 allocs/op
gc 49 @1.499s 0%: 0.019+0.044+0.020 ms clock, 0.076+0/0.037/0.054+0.083 ms cpu, 1->1->0 MB, 4 MB goal, 4 P (forced)
BenchmarkReverseEmptyString-4   	gc 50 @1.500s 0%: 0.045+0.044+0.008 ms clock, 0.18+0/0.034/0.054+0.035 ms cpu, 0->0->0 MB, 4 MB goal, 4 P (forced)
gc 51 @1.500s 0%: 0.004+0.030+0.009 ms clock, 0.017+0/0.024/0.053+0.037 ms cpu, 0->0->0 MB, 4 MB goal, 4 P (forced)
gc 52 @1.500s 0%: 0.001+0.036+0.009 ms clock, 0.005+0/0.023/0.054+0.037 ms cpu, 0->0->0 MB, 4 MB goal, 4 P (forced)
gc 53 @1.529s 0%: 0.018+0.040+0.020 ms clock, 0.074+0/0.035/0.058+0.080 ms cpu, 3->3->0 MB, 4 MB goal, 4 P (forced)
gc 54 @1.567s 0%: 0.003+1.0+0.043 ms clock, 0.012+0.046/0.020/0.033+0.17 ms cpu, 4->4->0 MB, 5 MB goal, 4 P
<snip>
gc 92 @2.955s 0%: 0.002+0.075+0.041 ms clock, 0.009+0.044/0.022/0.035+0.16 ms cpu, 4->4->0 MB, 5 MB goal, 4 P
50000000	        28.8 ns/op	       3 B/op	       1 allocs/op
gc 93 @2.969s 0%: 0.034+0.043+0.020 ms clock, 0.13+0/0.035/0.050+0.083 ms cpu, 1->1->0 MB, 4 MB goal, 4 P (forced)
BenchmarkReverseRunes-4         	gc 94 @2.970s 0%: 0.022+0.047+0.030 ms clock, 0.090+0/0.039/0.048+0.12 ms cpu, 0->0->0 MB, 4 MB goal, 4 P (forced)
gc 95 @2.970s 0%: 0.001+0.034+0.009 ms clock, 0.005+0/0.027/0.056+0.037 ms cpu, 0->0->0 MB, 4 MB goal, 4 P (forced)
gc 96 @2.971s 0%: 0.004+0.49+0.011 ms clock, 0.019+0/0.039/0.035+0.045 ms cpu, 0->0->0 MB, 4 MB goal, 4 P (forced)
gc 97 @3.006s 0%: 0.002+0.35+0.026 ms clock, 0.008+0.045/0.023/0.032+0.10 ms cpu, 4->4->0 MB, 5 MB goal, 4 P
gc 98 @3.041s 0%: 0.002+0.069+0.041 ms clock, 0.009+0.045/0.024/0.030+0.16 ms cpu, 4->4->0 MB, 5 MB goal, 4 P
gc 99 @3.075s 0%: 0.002+1.5+0.022 ms clock, 0.009+0.061/0.022/0.013+0.088 ms cpu, 4->4->0 MB, 5 MB goal, 4 P
gc 100 @3.108s 0%: 0.019+0.048+0.021 ms clock, 0.076+0/0.043/0.059+0.086 ms cpu, 3->3->0 MB, 4 MB goal, 4 P (forced)
gc 101 @3.143s 0%: 0.002+0.076+0.041 ms clock, 0.009+0.045/0.024/0.033+0.16 ms cpu, 4->4->0 MB, 5 MB goal, 4 P
<snip>
gc 139 @4.454s 0%: 0.002+0.064+0.016 ms clock, 0.009+0.051/0.020/0.029+0.066 ms cpu, 4->4->0 MB, 5 MB goal, 4 P
10000000	       135 ns/op	      16 B/op	       1 allocs/op
PASS
coverage: 100.0% of statements
ok  	_/home/owner/Documents/GO/TDD_and_Benchmarking_in_Go/stringutil	4.469s
$
```

[The above results were condensed slightly. Only non-forced GCs were elided.]

By my count, GC ran 152 times during this 4.5 second test (which kind of illustrates why GC must
be taken into account). We can get even more granular by setting __scavenge=1__ but this is kind
of overkill here (still, probably useful if checking for memory leaks).

Maybe we can turn Garbage Collection off?

```bash
$ GODEBUG=gctrace=1 GOGC=-1 go test -v -bench=. -cover -benchmem
=== RUN   TestReverse
--- PASS: TestReverse (0.00s)
gc 1 @0.000s 12%: 0.042+0.042+0.008 ms clock, 0.17+0/0.034/0.072+0.034 ms cpu, 0->0->0 MB, 17592186044415 MB goal, 4 P (forced)
gc 2 @0.000s 12%: 0.001+0.032+0.008 ms clock, 0.006+0/0.029/0.053+0.034 ms cpu, 0->0->0 MB, 17592186044415 MB goal, 4 P (forced)
goos: linux
goarch: amd64
BenchmarkReverseBytes-4         	gc 3 @0.000s 12%: 0.001+0.034+0.009 ms clock, 0.005+0/0.024/0.054+0.038 ms cpu, 0->0->0 MB, 17592186044415 MB goal, 4 P (forced)
gc 4 @0.000s 11%: 0.001+0.032+0.008 ms clock, 0.006+0/0.025/0.051+0.033 ms cpu, 0->0->0 MB, 17592186044415 MB goal, 4 P (forced)
gc 5 @0.002s 5%: 0.001+0.035+0.008 ms clock, 0.006+0/0.029/0.059+0.034 ms cpu, 0->0->0 MB, 17592186044415 MB goal, 4 P (forced)
gc 6 @0.141s 0%: 0.019+0.067+0.017 ms clock, 0.079+0/0.061/0.091+0.069 ms cpu, 15->15->0 MB, 17592186044415 MB goal, 4 P (forced)
10000000	       138 ns/op	      16 B/op	       1 allocs/op
gc 7 @1.524s 0%: 0.010+0.30+0.043 ms clock, 0.043+0/0.29/0.54+0.17 ms cpu, 152->152->0 MB, 17592186044415 MB goal, 4 P (forced)
BenchmarkReverseEmptyString-4   	gc 8 @1.532s 0%: 0.021+0.046+0.008 ms clock, 0.087+0/0.037/0.049+0.035 ms cpu, 0->0->0 MB, 17592186044415 MB goal, 4 P (forced)
gc 9 @1.532s 0%: 0.002+0.040+0.012 ms clock, 0.010+0/0.030/0.047+0.049 ms cpu, 0->0->0 MB, 17592186044415 MB goal, 4 P (forced)
gc 10 @1.532s 0%: 0.001+0.030+0.008 ms clock, 0.006+0/0.027/0.049+0.033 ms cpu, 0->0->0 MB, 17592186044415 MB goal, 4 P (forced)
gc 11 @1.561s 0%: 0.020+0.051+0.030 ms clock, 0.081+0/0.042/0.059+0.12 ms cpu, 3->3->0 MB, 17592186044415 MB goal, 4 P (forced)
50000000	        28.5 ns/op	       3 B/op	       1 allocs/op
gc 12 @2.987s 0%: 0.018+0.30+0.029 ms clock, 0.073+0/0.29/0.54+0.11 ms cpu, 152->152->0 MB, 17592186044415 MB goal, 4 P (forced)
BenchmarkReverseRunes-4         	gc 13 @2.995s 0%: 0.043+0.045+0.023 ms clock, 0.17+0/0.039/0.050+0.093 ms cpu, 0->0->0 MB, 17592186044415 MB goal, 4 P (forced)
gc 14 @2.995s 0%: 0.003+0.040+0.009 ms clock, 0.012+0/0.034/0.046+0.036 ms cpu, 0->0->0 MB, 17592186044415 MB goal, 4 P (forced)
gc 15 @2.997s 0%: 0.002+0.032+0.009 ms clock, 0.008+0/0.027/0.049+0.036 ms cpu, 0->0->0 MB, 17592186044415 MB goal, 4 P (forced)
gc 16 @3.132s 0%: 0.019+0.064+0.043 ms clock, 0.076+0/0.054/0.091+0.17 ms cpu, 15->15->0 MB, 17592186044415 MB goal, 4 P (forced)
10000000	       135 ns/op	      16 B/op	       1 allocs/op
PASS
coverage: 100.0% of statements
ok  	_/home/owner/Documents/GO/TDD_and_Benchmarking_in_Go/stringutil	4.490s
$
```

Okay, this eliminates all of the non-forced GCs (which is something - reducing GCs from 152 to 16).

Which is not exactly what the documentation says:

> A negative percentage disables garbage collection.

From:

    http://golang.org/pkg/runtime/debug/#SetGCPercent

Still, eliminating all non-forced GCs is something. It had a minimal impact on run time (4.469s to 4.490s),
which evaluates to about 0.47 percent (this may simply be the time required to report on GC).

So about half a percent.

#### Alternatives

Okay, now that we've established our benchmarking procedure, lets verify things with some
alternative implementations, which we will call __naive\_reverse.go__ and __better\_reverse.go__.
These alternative implementations will have tests and benchmarks, which will all be the same.

We are not expecting to improve upon the original code (which was written by Google, after
all), merely to have some alternatives to benchmark.

And the results are:

```bash
$ GOGC=-1 go test -v -bench=. -cover -benchmem
=== RUN   TestBetterReverse
--- PASS: TestBetterReverse (0.00s)
=== RUN   TestNaiveReverse
--- PASS: TestNaiveReverse (0.00s)
=== RUN   TestReverse
--- PASS: TestReverse (0.00s)
goos: linux
goarch: amd64
BenchmarkBetterReverseBytes-4         	10000000	       173 ns/op	      64 B/op	       2 allocs/op
BenchmarkBetterReverseEmptyString-4   	50000000	        34.6 ns/op	       3 B/op	       1 allocs/op
BenchmarkBetterReverseRunes-4         	10000000	       172 ns/op	      64 B/op	       2 allocs/op
BenchmarkNaiveReverseBytes-4          	10000000	       163 ns/op	      64 B/op	       2 allocs/op
BenchmarkNaiveReverseEmptyString-4    	50000000	        34.2 ns/op	       3 B/op	       1 allocs/op
BenchmarkNaiveReverseRunes-4          	10000000	       164 ns/op	      64 B/op	       2 allocs/op
BenchmarkReverseBytes-4               	10000000	       141 ns/op	      16 B/op	       1 allocs/op
BenchmarkReverseEmptyString-4         	50000000	        28.5 ns/op	       3 B/op	       1 allocs/op
BenchmarkReverseRunes-4               	10000000	       135 ns/op	      16 B/op	       1 allocs/op
PASS
coverage: 100.0% of statements
ok  	_/home/owner/Documents/GO/TDD_and_Benchmarking_in_Go/stringutil	15.632s
$
```

And it seems that our so-called ___naive___ implementation is actually better than our
so-called ___better___ implementation.

[It is for this reason that people say "premature optimization is the root of all evil".
It is all too easy to "optimize" code that is not actually a bottleneck; likewise it is
equally easy to make "improvements" to code that actually makes it slower.]

As expected, neither of our new implementations comes close to the original.

The key reason for this is that the original uses an ___atomic___ operation to swap
runes in-place - which saves a second memory allocation. As well, this pretty much
halves the number of iterations - for a small speed increase.

Here is the particularly clever part of the code:

```Golang
		r[i], r[j] = r[j], r[i]
```

Just for fun, let's try this again with a more recent version of Go (1.12.5):

```bash
$ GOGC=-1 go test -v -bench=. -cover -benchmem
=== RUN   TestBetterReverse
--- PASS: TestBetterReverse (0.00s)
=== RUN   TestNaiveReverse
--- PASS: TestNaiveReverse (0.00s)
=== RUN   TestReverse
--- PASS: TestReverse (0.00s)
goos: linux
goarch: amd64
BenchmarkBetterReverseBytes-4         	10000000	       170 ns/op	      64 B/op	       2 allocs/op
BenchmarkBetterReverseEmptyString-4   	50000000	        34.1 ns/op	       3 B/op	       1 allocs/op
BenchmarkBetterReverseRunes-4         	10000000	       179 ns/op	      64 B/op	       2 allocs/op
BenchmarkNaiveReverseBytes-4          	10000000	       160 ns/op	      64 B/op	       2 allocs/op
BenchmarkNaiveReverseEmptyString-4    	50000000	        35.4 ns/op	       3 B/op	       1 allocs/op
BenchmarkNaiveReverseRunes-4          	10000000	       167 ns/op	      64 B/op	       2 allocs/op
BenchmarkReverseBytes-4               	10000000	       133 ns/op	      16 B/op	       1 allocs/op
BenchmarkReverseEmptyString-4         	50000000	        31.1 ns/op	       3 B/op	       1 allocs/op
BenchmarkReverseRunes-4               	10000000	       142 ns/op	      16 B/op	       1 allocs/op
PASS
coverage: 100.0% of statements
ok  	_/home/owner/Documents/GO/TDD_and_Benchmarking_in_Go/stringutil	15.851s
$
```

And there are some minor differences, but the overall results hold. Garbage collection in Go
does get tweaked from time to time, however the overall execution time is about the same so
it seems it hasn't been changed much with recent versions of Go.

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

## Versions

The version of Golang used for this exercise is __1.11__.

Verify the version of Go as follows:

```bash
$ go version
go version go1.11 linux/amd64
$
```

Different versions of Go can be expected to have different characteristics with
respect to how they handle Garbage Collection.

## To Do

- [x] Implement benchmarks
- [x] Create different implementations of the code & benchmark them
- [x] Re-run benchmarks with a more recent version of Go (1.12.5)
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
 desirable or even possible. Nevertheless, more code coverage is
 almost always (but not always) better than less code coverage.]
