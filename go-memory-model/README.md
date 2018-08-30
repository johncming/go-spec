...

<!-- #lowframe -->

[The Go Programming Language](https://golang.org/)

[Go](https://golang.org/)[▽](https://golang.org/ref/mem#)

<form method="GET" action="/search">
    [Documents](https://golang.org/doc/)
    [Packages](https://golang.org/pkg/)
    [The Project](https://golang.org/project/)
    [Help](https://golang.org/help/)

    [Blog](https://golang.org/blog/)


    [Play](http://play.golang.org/ "Show Go Playground")

    <input type="search" id="search" name="q" placeholder="Search" aria-label="Search" required="">

    <button type="submit"><span><!-- magnifying glass: --><svg width="24" height="24" viewBox="0 0 24 24"><title>submit search</title><path d="M15.5 14h-.79l-.28-.27C15.41 12.59 16 11.11 16 9.5 16 5.91 13.09 3 9.5 3S3 5.91 3 9.5 5.91 16 9.5 16c1.61 0 3.09-.59 4.23-1.57l.27.28v.79l5 4.99L20.49 19l-4.99-5zm-6 0C7.01 14 5 11.99 5 9.5S7.01 5 9.5 5 14 7.01 14 9.5 11.99 14 9.5 14z"></path><path d="M0 0h24v24H0z" fill="none"></path></svg></span></button>
</form>

<textarea class="code" spellcheck="false">package main

import "fmt"

func main() {
	fmt.Println("Hello, 世界")
}</textarea>

Run
		Format

		Share

#     The Go Memory Model

## Version of May 31, 2014

col 1                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                  | col 2
-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------
<dl>
    <dt>[Introduction](https://golang.org/ref/mem#tmp_0)</dt>
    <dt>[Advice](https://golang.org/ref/mem#tmp_1)</dt>
    <dt>[Happens Before](https://golang.org/ref/mem#tmp_2)</dt>
    <dt>[Synchronization](https://golang.org/ref/mem#tmp_3)</dt>
    <dd class="indent">[Initialization](https://golang.org/ref/mem#tmp_4)</dd>
    <dd class="indent">[Goroutine creation](https://golang.org/ref/mem#tmp_5)</dd>
    <dd class="indent">[Goroutine destruction](https://golang.org/ref/mem#tmp_6)</dd>
    <dd class="indent">[Channel communication](https://golang.org/ref/mem#tmp_7)</dd>
    <dd class="indent">[Locks](https://golang.org/ref/mem#tmp_8)</dd>
    <dd class="indent">[Once](https://golang.org/ref/mem#tmp_9)</dd>
    <dt>[Incorrect synchronization](https://golang.org/ref/mem#tmp_10)</dt>
</dl> | <dl></dl>

## Introduction

The Go memory model specifies the conditions under which
reads of a variable in one goroutine can be guaranteed to
observe values produced by writes to the same variable in a different goroutine.

## Advice

Programs that modify data being simultaneously accessed by multiple goroutines
must serialize such access.

To serialize access, protect the data with channel operations or other synchronization primitives
such as those in the [`sync`](https://golang.org/pkg/sync/)
and [`sync/atomic`](https://golang.org/pkg/sync/atomic/) packages.

If you must read the rest of this document to understand the behavior of your program,
you are being too clever.

Don't be clever.

## Happens Before

Within a single goroutine, reads and writes must behave
as if they executed in the order specified by the program.
That is, compilers and processors may reorder the reads and writes
executed within a single goroutine only when the reordering
does not change the behavior within that goroutine
as defined by the language specification.
Because of this reordering, the execution order observed
by one goroutine may differ from the order perceived
by another.  For example, if one goroutine
executes `a = 1; b = 2;`, another might observe
the updated value of `b` before the updated value of `a`.

To specify the requirements of reads and writes, we define
_happens before_, a partial order on the execution
of memory operations in a Go program.  If event e<sub>1</sub> happens
before event e<sub>2</sub>, then we say that e<sub>2</sub> happens after e<sub>1</sub>.
Also, if e<sub>1</sub> does not happen before e<sub>2</sub> and does not happen
after e<sub>2</sub>, then we say that e<sub>1</sub> and e<sub>2</sub> happen concurrently.

Within a single goroutine, the happens-before order is the
order expressed by the program.

A read r of a variable `v` is _allowed_ to observe a write w to `v`
if both of the following hold:

1.  r does not happen before w.
2.  There is no other write w' to `v` that happens
        after w but before r.

To guarantee that a read r of a variable `v` observes a
particular write w to `v`, ensure that w is the only
write r is allowed to observe.
That is, r is _guaranteed_ to observe w if both of the following hold:

1.  w happens before r.
2.  Any other write to the shared variable `v`
    either happens before w or after r.

This pair of conditions is stronger than the first pair;
it requires that there are no other writes happening
concurrently with w or r.

Within a single goroutine,
there is no concurrency, so the two definitions are equivalent:
a read r observes the value written by the most recent write w to `v`.
When multiple goroutines access a shared variable `v`,
they must use synchronization events to establish
happens-before conditions that ensure reads observe the
desired writes.

The initialization of variable `v` with the zero value
for `v`'s type behaves as a write in the memory model.

Reads and writes of values larger than a single machine word
behave as multiple machine-word-sized operations in an
unspecified order.

## Synchronization

### Initialization

Program initialization runs in a single goroutine,
but that goroutine may create other goroutines,
which run concurrently.

If a package `p` imports package `q`, the completion of
`q`'s `init` functions happens before the start of any of `p`'s.

The start of the function `main.main` happens after
all `init` functions have finished.

### Goroutine creation

The `go` statement that starts a new goroutine
happens before the goroutine's execution begins.

For example, in this program:

<pre>
var a string

func f() {
	print(a)
}

func hello() {
	a = "hello, world"
	go f()
}
</pre>

calling `hello` will print `"hello, world"`
at some point in the future (perhaps after `hello` has returned).

### Goroutine destruction

The exit of a goroutine is not guaranteed to happen before
any event in the program.  For example, in this program:

<pre>
var a string

func hello() {
	go func() { a = "hello" }()
	print(a)
}
</pre>

the assignment to `a` is not followed by
any synchronization event, so it is not guaranteed to be
observed by any other goroutine.
In fact, an aggressive compiler might delete the entire `go` statement.

If the effects of a goroutine must be observed by another goroutine,
use a synchronization mechanism such as a lock or channel
communication to establish a relative ordering.

### Channel communication

Channel communication is the main method of synchronization
between goroutines.  Each send on a particular channel
is matched to a corresponding receive from that channel,
usually in a different goroutine.

A send on a channel happens before the corresponding
receive from that channel completes.

This program:

<pre>
var c = make(chan int, 10)
var a string

func f() {
	a = "hello, world"
	c <- 0
}

func main() {
	go f()
	<-c
	print(a)
}
</pre>

is guaranteed to print `"hello, world"`.  The write to `a`
happens before the send on `c`, which happens before
the corresponding receive on `c` completes, which happens before
the `print`.

The closing of a channel happens before a receive that returns a zero value
because the channel is closed.

In the previous example, replacing
`c <- 0` with `close(c)`
yields a program with the same guaranteed behavior.

A receive from an unbuffered channel happens before
the send on that channel completes.

This program (as above, but with the send and receive statements swapped and
using an unbuffered channel):

<pre>
var c = make(chan int)
var a string

func f() {
	a = "hello, world"
	<-c
}
</pre>

<pre>
func main() {
	go f()
	c <- 0
	print(a)
}
</pre>

is also guaranteed to print `"hello, world"`.  The write to `a`
happens before the receive on `c`, which happens before
the corresponding send on `c` completes, which happens
before the `print`.

If the channel were buffered (e.g., `c = make(chan int, 1)`)
then the program would not be guaranteed to print
`"hello, world"`.  (It might print the empty string,
crash, or do something else.)

The _k_th receive on a channel with capacity _C_ happens before the _k_+_C_th send from that channel completes.

This rule generalizes the previous rule to buffered channels.
It allows a counting semaphore to be modeled by a buffered channel:
the number of items in the channel corresponds to the number of active uses,
the capacity of the channel corresponds to the maximum number of simultaneous uses,
sending an item acquires the semaphore, and receiving an item releases
the semaphore.
This is a common idiom for limiting concurrency.

This program starts a goroutine for every entry in the work list, but the
goroutines coordinate using the `limit` channel to ensure
that at most three are running work functions at a time.

<pre>
var limit = make(chan int, 3)

func main() {
	for _, w := range work {
		go func(w func()) {
			limit <- 1
			w()
			<-limit
		}(w)
	}
	select{}
}
</pre>

### Locks

The `sync` package implements two lock data types,
`sync.Mutex` and `sync.RWMutex`.

For any `sync.Mutex` or `sync.RWMutex` variable `l` and _n_ < _m_,
call _n_ of `l.Unlock()` happens before call _m_ of `l.Lock()` returns.

This program:

<pre>
var l sync.Mutex
var a string

func f() {
	a = "hello, world"
	l.Unlock()
}

func main() {
	l.Lock()
	go f()
	l.Lock()
	print(a)
}
</pre>

is guaranteed to print `"hello, world"`.
The first call to `l.Unlock()` (in `f`) happens
before the second call to `l.Lock()` (in `main`) returns,
which happens before the `print`.

For any call to `l.RLock` on a `sync.RWMutex` variable `l`,
there is an _n_ such that the `l.RLock` happens (returns) after call _n_ to
`l.Unlock` and the matching `l.RUnlock` happens
before call _n_+1 to `l.Lock`.

### Once

The `sync` package provides a safe mechanism for
initialization in the presence of multiple goroutines
through the use of the `Once` type.
Multiple threads can execute `once.Do(f)` for a particular `f`,
but only one will run `f()`, and the other calls block
until `f()` has returned.

A single call of `f()` from `once.Do(f)` happens (returns) before any call of `once.Do(f)` returns.

In this program:

<pre>
var a string
var once sync.Once

func setup() {
	a = "hello, world"
}

func doprint() {
	once.Do(setup)
	print(a)
}

func twoprint() {
	go doprint()
	go doprint()
}
</pre>

calling `twoprint` causes `"hello, world"` to be printed twice.
The first call to `doprint` runs `setup` once.

## Incorrect synchronization

Note that a read r may observe the value written by a write w
that happens concurrently with r.
Even if this occurs, it does not imply that reads happening after r
will observe writes that happened before w.

In this program:

<pre>
var a, b int

func f() {
	a = 1
	b = 2
}

func g() {
	print(b)
	print(a)
}

func main() {
	go f()
	g()
}
</pre>

it can happen that `g` prints `2` and then `0`.

This fact invalidates a few common idioms.

Double-checked locking is an attempt to avoid the overhead of synchronization.
For example, the `twoprint` program might be
incorrectly written as:

<pre>
var a string
var done bool

func setup() {
	a = "hello, world"
	done = true
}

func doprint() {
	if !done {
		once.Do(setup)
	}
	print(a)
}

func twoprint() {
	go doprint()
	go doprint()
}
</pre>

but there is no guarantee that, in `doprint`, observing the write to `done`
implies observing the write to `a`.  This
version can (incorrectly) print an empty string
instead of `"hello, world"`.

Another incorrect idiom is busy waiting for a value, as in:

<pre>
var a string
var done bool

func setup() {
	a = "hello, world"
	done = true
}

func main() {
	go setup()
	for !done {
	}
	print(a)
}
</pre>

As before, there is no guarantee that, in `main`,
observing the write to `done`
implies observing the write to `a`, so this program could
print an empty string too.
Worse, there is no guarantee that the write to `done` will ever
be observed by `main`, since there are no synchronization
events between the two threads.  The loop in `main` is not
guaranteed to finish.

There are subtler variants on this theme, such as this program.

<pre>
type T struct {
	msg string
}

var g *T

func setup() {
	t := new(T)
	t.msg = "hello, world"
	g = t
}

func main() {
	go setup()
	for g == nil {
	}
	print(g.msg)
}
</pre>

Even if `main` observes `g != nil` and exits its loop,
there is no guarantee that it will observe the initialized
value for `g.msg`.

In all these examples, the solution is the same:
use explicit synchronization.

Build version go1.11.

Except as [noted](https://developers.google.com/site-policies#restrictions),
the content of this page is licensed under the
Creative Commons Attribution 3.0 License,
and code is licensed under a [BSD license](https://golang.org/LICENSE).
[Terms of Service](https://golang.org/doc/tos.html) |
[Privacy Policy](http://www.google.com/intl/en/policies/privacy/)
