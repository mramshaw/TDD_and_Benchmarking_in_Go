# TDD and Benchmarking in Go

I was experimenting with [TinyGo](http://github.com/mramshaw/TinyGo) and thought I'd have a go at
string reversal. This is pretty easy to do __bytewise__ but not yet possible __runewise__.

Based upon recent discussions about TDD & Testing and looking at the example code in:

    http://github.com/golang/example/tree/master/stringutil

It seemed to me that the tests implemented here do not meet the criteria of what *I* would do.

Accordingly, it seemed to be time to have a quick stab at TDD and Benchmarking in Golang, using
these particular pieces of code as a basis.

The stated goals for this code are:

* The basic form of a library
* Conversion between string and []rune
* Table-driven unit tests ([testing](http://golang.org/pkg/testing/))

My intent here is to examine (and possibly improve) the third goal.

## To Do

- [ ] Implement __runes__ in TinyGo

## Credits

Somewhat inspired by this fascinating podcast ("Compilers, LLVM, Swift, TPU, and ML Accelerators"):

    http://lexfridman.com/chris-lattner/

[Lex Fridman usually has fascinating guests & really interesting questions.]

The key reminder was that it is possible to optimize for different use cases.

Also by this interesting podcast ("It's time to talk about testing"):

    http://changelog.com/gotime/83

[The inspiration here was the discussion around whether or not 100% code
 coverage was really sufficient. TL;DR - it really depends upon the test
 quality.]
