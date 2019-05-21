# TDD and Benchmarking in Go

I was experimenting with [TinyGo](http://github.com/mramshaw/TinyGo) and thought I'd take a stab
at string reversal. This is pretty easy to do __bytewise__ in TinyGo but not yet possible __runewise__
(rune conversion is on the TinyGo todo list).

Based upon recent discussions about TDD & Testing and looking at the example code in:

    http://github.com/golang/example/tree/master/stringutil

It seemed to me that the tests implemented there do not meet the criteria of what ***I*** would do.

Accordingly, it seemed to be time for a refresher on TDD and Benchmarking in Golang using
this code as a basis.

The stated goals for the code are:

* The basic form of a library
* Conversion between string and []rune
* Table-driven unit tests ([testing](http://golang.org/pkg/testing/))

My intent here is to examine (and possibly improve) the third goal.

## Discussion

Some points to keep in mind are as follows.

[These are from the __builtin__ package documentation, which is where native Golang data
types - which would otherwise be undocumented - are documented.]

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

- [ ] Implement better tests
- [ ] Create different implementations of the code & benchmark them
- [ ] Verify the optimization criteria for TinyGo
- [ ] Implement __runes__ in TinyGo

## Credits

Somewhat inspired by this fascinating podcast ("Compilers, LLVM, Swift, TPU, and ML Accelerators"):

    http://lexfridman.com/chris-lattner/

[Lex Fridman usually has fascinating guests & prepares really interesting questions.
 The key reminder was that it is possible to optimize for different use cases.]

Also by this interesting podcast ("It's time to talk about testing"):

    http://changelog.com/gotime/83

[The inspiration here was the discussion around whether or not 100% code
 coverage was really sufficient. TL;DR - it really depends upon the test
 quality.]
