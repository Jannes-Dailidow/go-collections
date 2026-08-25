# go-collections

Generic collection types with a shared iterator layer, built on Go's `iter` package.

```sh
go get github.com/jannes-dailidow/go-collections
```

```go
import "github.com/jannes-dailidow/go-collections"
```

Requires **Go 1.27**, for generic methods.

```go
users := collections.Slice[User]{ /* ... */ }

// lazy all the way to the Slice() at the end
names := users.Values().
	Filter(func(u User) bool { return u.Active }).
	SortBy(func(u User) string { return u.Surname }).
	Map(func(u User) string { return u.Name }).
	Take(10).
	Slice()

// a fallible transform moves into the X family and carries one error out
ids, err := users.Values().MapX(parseID).Slice()

// group, keeping the groups in the order their keys first appeared
byTeam := collections.GroupByOrdered(users.Values(), func(u User) string { return u.Team })

// the collections answer the O(1) questions an iterator cannot answer cheaply
teams := collections.CollectOrderedSet(byTeam.Keys())
teams.Has("platform")     // O(1)
teams.IndexOf("platform") // O(1), through the internal index
```

Runnable versions of these live in `example_test.go`, so they are compiled and their
output checked with the rest of the suite.

## Contents

- [Layers](#layers) — how the package is organised
- [Conventions](#conventions) — the naming rules, and the `X` suffix
- [Types](#types) — the nine types and their entry points
- [Operations](#operations) — the full surface, by category
- [Technical decisions](#technical-decisions) — why it is shaped this way
- [Migrating from v0.1](#migrating-from-v01)
- [Testing](#testing)
- [Project layout](#project-layout)
- [Known gaps](#known-gaps)

## Layers

Two layers, and the split is the whole design.

**Five collections** — `Slice`, `Map`, `Set`, `OrderedMap`, `OrderedSet` — carry only
what they can do better than a generic iterator: O(1) lookup, set algebra, positional
access, in-place mutation. Plus the entry points into the other layer.

**Four iterator types** — `Iter`, `Iter2`, `IterX`, `Iter2X` — carry every operation.
`Values()`, `All()` and `Keys()` are one call away from any collection, so a
collection does not re-expose `Map`, `Filter` and the rest. Duplicating sixty
operations across five collections would be five times the code for no new power.

```text
Slice[T]  Map[K,V]  Set[T]  OrderedMap[K,V]  OrderedSet[T]
    |         |        |            |              |
    +---------+--------+------------+--------------+
                       |  Values() / All() / Keys()
                       v
          Iter[T]  <->  Iter2[K,V]        infallible
             ^              ^
             | X()          | X()
             v              v
          IterX[T] <-> Iter2X[K,V]        can abort with an error
```

## Conventions

| Name | Meaning |
| ---- | ------- |
| `Native()` | the plain Go type — `[]T`, `map[K]V`, `map[T]struct{}` |
| `Seq()` / `Seq2()` | the stdlib `iter.Seq` / `iter.Seq2` equivalent |
| `All()` / `Keys()` / `Values()` | this package's `Iter` / `Iter2` |
| `Collect*` | builds a collection by draining an iterator |
| `*By` | takes a key function, so a constraint can land on the key |
| `X` | fallible — see below |

### The X suffix means two things

On a **type** it means the stream can abort: `IterX[T]` is
`func(yield func(T) bool) error`.

On a **method** it means the callback can fail, which lifts an infallible stream into
the fallible family:

| Receiver | Method | Result |
| -------- | ------ | ------ |
| `Iter[T]` | `Map(fn func(T) U)` | `Iter[U]` |
| `Iter[T]` | `MapX(fn func(T) (U, error))` | `IterX[U]` |
| `IterX[T]` | `Map(fn func(T) U)` | `IterX[U]` |
| `IterX[T]` | `MapX(fn func(T) (U, error))` | `IterX[U]` |

An abort is **not a value in the stream**. The producer returns it to stop iteration,
it surfaces once at the end, and collectors discard the partial result. A consumer
`break` is not an abort and leaves the error nil.

An `X` iterator is not directly rangeable. `Iter(&err)` captures the abort
out-of-band so that it can be:

```go
var err error
for v := range stream.Iter(&err) {
    // ...
}
if err != nil { /* ... */ }
```

A nil pointer discards the error.

## Types

### Slice[T]

A `[]T` that can enter an iterator pipeline without conversion.

| Method | Returns |
| ------ | ------- |
| `Native()` | `[]T` |
| `Values()` | `Iter[T]` |
| `All()` | `Iter2[int, T]` |
| `Backward()` | `Iter2[int, T]`, reverse order, no copy |
| `Seq()` / `Seq2()` | `iter.Seq[T]` / `iter.Seq2[int, T]` |
| `Len()` / `IsEmpty()` | `int` / `bool`, O(1) |
| `At(n)` | `*T`, nil out of range, pointing **into** the slice |
| `Append(v...)` / `Insert(n, v...)` / `DeleteAt(n)` | `Slice[T]` |
| `Clone()` / `Grow(n)` | `Slice[T]` |
| `SortBy(fn)` / `SortFunc(cmp)` / `Reverse()` | `Slice[T]`, **in place** |
| `SetBy(fn)` / `MapBy(fn)` | `Set[K]` / `Map[K, V]` |
| `ParallelMap` / `ParallelFilter` / `ParallelEach` | see [Parallel](#parallel) |

### Map[K, V]

A `map[K]V` with key, value and pair iteration. Iteration order is random.

| Method | Returns |
| ------ | ------- |
| `CollectMap(Iter2[K, V])` / `CollectMapX(Iter2X[K, V])` | `Map[K, V]` / `(Map[K, V], error)` |
| `Native()` | `map[K]V` |
| `Keys()` / `Values()` | `Iter[K]` / `Iter[V]` |
| `All()` / `Seq2()` | `Iter2[K, V]` / `iter.Seq2[K, V]` |
| `OrderedMap()` | `*OrderedMap[K, V]`, in arbitrary order |
| `Len()` / `IsEmpty()` / `Has(k)` | O(1) |
| `Get(k)` / `GetOr(k, fallback)` | `(V, bool)` / `V` |
| `Put(k, v)` / `Delete(k)` / `Clear()` | — |
| `Clone()` | `Map[K, V]`, shallow |
| `Merge(other)` / `MergeFunc(other, resolve)` | — |
| `InvertBy(fn)` | `Map[V2, K]` |

`Put` panics on a nil map, like the assignment it wraps. Build one with
`make(Map[K, V])` or a `Collect` function.

### Set[T]

Membership without values, as a zero-width `map[T]struct{}`. Iteration order is
random.

| Method | Returns |
| ------ | ------- |
| `CollectSet(Iter[T])` / `CollectSetX(IterX[T])` | `Set[T]` / `(Set[T], error)` |
| `Native()` / `Values()` / `Slice()` / `Seq()` | `map[T]struct{}` / `Iter[T]` / `Slice[T]` / `iter.Seq[T]` |
| `OrderedSet()` | `*OrderedSet[T]`, in arbitrary order |
| `Len()` / `IsEmpty()` / `Has(v)` | O(1) |
| `Add(v)` / `Remove(v)` | `bool`, reporting whether it changed |
| `Clear()` / `Clone()` | — / `Set[T]` |
| `Union` / `Intersect` / `Diff` / `SymDiff` | `Set[T]` |
| `IsSubset` / `IsSuperset` / `IsDisjoint` | `bool` |

### OrderedMap[K, V]

A map that iterates in insertion order. Re-inserting a key updates its value and
keeps its original position. Every read is safe on a nil receiver.

| Method | Returns |
| ------ | ------- |
| `NewOrderedMap[K, V]()` | `*OrderedMap[K, V]`; the zero value works too |
| `CollectOrderedMap(Iter2[K, V])` / `...X` | `*OrderedMap[K, V]` / `(..., error)` |
| `All()` / `Keys()` / `Values()` / `Seq2()` | as `Map` |
| `Backward()` | `Iter2[K, V]`, last insertion first |
| `Map()` / `Native()` | `Map[K, V]` / `map[K]V` |
| `Len()` / `IsEmpty()` / `Has(k)` / `IndexOf(k)` | O(1) |
| `Get(k)` / `GetOr(k, fallback)` | `(V, bool)` / `V` |
| `At(n)` | `*Entry[K, V]`, nil out of range, a **copy** |
| `Put(k, v)` / `Delete(k)` / `Clear()` | `Delete` is O(n) |
| `Clone()` | `*OrderedMap[K, V]` |
| `SortBy(fn)` / `SortFunc(cmp)` / `Reverse()` | `*OrderedMap[K, V]`, in place |

### OrderedSet[T]

A set that iterates in first-seen order. Every read is safe on a nil receiver.

| Method | Returns |
| ------ | ------- |
| `NewOrderedSet[T]()` | `*OrderedSet[T]`; the zero value works too |
| `CollectOrderedSet(Iter[T])` / `...X` | `*OrderedSet[T]` / `(..., error)` |
| `Values()` / `All()` / `Backward()` / `Seq()` | `Iter[T]` / `Iter2[int, T]` / `Iter2[int, T]` / `iter.Seq[T]` |
| `Slice()` / `Native()` / `Set()` | a copy, not a view of the internal order |
| `Len()` / `IsEmpty()` / `Has(v)` / `IndexOf(v)` | O(1) |
| `At(n)` | `*T`, nil out of range, a **copy** |
| `Add(v)` / `Remove(v)` | `bool`; `Remove` is O(n) |
| `Clear()` / `Clone()` | — / `*OrderedSet[T]` |
| `Union` / `Intersect` / `Diff` / `SymDiff` | `*OrderedSet[T]`, receiver's order first |
| `IsSubset` / `IsSuperset` / `IsDisjoint` | `bool` |
| `SortBy(fn)` / `SortFunc(cmp)` / `Reverse()` | `*OrderedSet[T]`, in place |

### Entry[K, V]

A key/value pair, with exported `Key` and `Value`. `OrderedMap` stores its order as
`[]Entry[K, V]`, and the pair operations hand it back from `Find`, `At`, `MinBy` and
`Pairs`.

### Iter[T] and Iter2[K, V]

`func(yield func(T) bool)` and `func(yield func(K, V) bool)`. Shaped like
`iter.Seq` / `iter.Seq2` and usable with `range`.

| Method | Returns |
| ------ | ------- |
| `FromSeq(iter.Seq[T])` / `FromSeq2(...)` | `Iter[T]` / `Iter2[K, V]` |
| `Seq()` / `Seq2()` | the stdlib equivalent |
| `X()` | `IterX[T]` / `Iter2X[K, V]`, where it can never abort |
| `Slice()` / `Native()` | `Slice[T]` / `[]T` (single-value only) |
| `Keys()` / `Values()` | `Iter[K]` / `Iter[V]` (pair only) |

### IterX[T] and Iter2X[K, V]

`func(yield func(T) bool) error` and `func(yield func(K, V) bool) error`.

| Method | Returns |
| ------ | ------- |
| `Iter(*error)` / `Iter2(*error)` | `Iter[T]` / `Iter2[K, V]`, rangeable |
| `Seq(*error)` / `Seq2(*error)` | the stdlib equivalent |
| `Slice()` / `Native()` | `(Slice[T], error)` / `([]T, error)` (single-value only) |
| `Keys()` / `Values()` | `IterX[K]` / `IterX[V]` — still fallible, so the abort survives |

## Operations

Every operation below is on all four iterator types unless the column says otherwise.
Collections reach them through `Values()`, `All()` or `Keys()`.

Shorthand: `Iter*` is all four, `Iter1*` is `Iter` and `IterX`, `Iter2*` is `Iter2`
and `Iter2X`.

### Transform

All lazy.

| Operation | On | Notes |
| --------- | -- | ----- |
| `Map` | `Iter*` | `Iter[U]`; on a pair stream it replaces both halves |
| `MapX` | `Iter*` | callback returns `(U, error)`, lifting into the `X` family |
| `MapKeys` / `MapValues` | `Iter2*` | replace one half; both callbacks see both halves |
| `KeyBy` | `Iter1*` | `Iter2[K, T]` — enter the pair layer |
| `Collapse` | `Iter2*` | `Iter[T]` — leave it |
| `Pairs(i)` | func over `Iter2*` | `Iter[Entry[K, V]]` — leave it keeping both halves |
| `Enumerate` | `Iter1*` | `Iter2[int, T]`, numbering the stream |
| `Zip(other)` | `Iter1*` | `Iter2[T, U]`, stopping at the shorter side |
| `FlatMap` | `Iter1*` | callback returns an iterator per element |
| `Flatten(i)` / `FlattenSlices(i)` | funcs | concatenate a stream of streams or of slices |
| `Scan` | `Iter1*` | the running accumulator, one value per element |
| `Chunk(i, n)` / `Window(i, n)` | funcs over `Iter1*` | `Iter[Slice[T]]`, partitioning or sliding |
| `Tap` | `Iter*` | a side effect per element, passed through |

### Select

All lazy, all shape-preserving.

| Operation | On | Notes |
| --------- | -- | ----- |
| `Filter` / `FilterX` | `Iter*` | |
| `Take(n)` / `Drop(n)` | `Iter*` | `Drop(10).Take(5)` paginates without materializing |
| `TakeWhile` / `DropWhile` | `Iter*` | with `X` forms |
| `DedupBy` | `Iter*` | first occurrence wins |
| `Dedup(i)` / `Compact(i)` | funcs over `Iter1*` | need `T comparable`; `CompactBy` is the method form |

### Reorder and group

These buffer the whole stream, so they will not terminate on an endless one.

| Operation | On | Notes |
| --------- | -- | ----- |
| `SortBy` / `SortFunc` | `Iter*`, `Slice`, `OrderedMap`, `OrderedSet` | stable; iterators buffer a copy, collections sort in place |
| `Reverse` | `Iter*`, `Slice`, `OrderedSet`, `OrderedMap` | on a `Slice`, `Backward()` is cheaper |
| `GroupBy(i, fn)` / `GroupByOrdered(i, fn)` | funcs over `Iter1*` | `Map[K, Slice[T]]` / `*OrderedMap[K, Slice[T]]` |
| `Partition` | `Iter1*` | `(accepted, rejected Slice[T])`; an error discards both |

### Terminal queries

These drain the stream and return a value. On the `X` types each gains a second
return.

| Operation | On | Notes |
| --------- | -- | ----- |
| `Each` / `EachX` | `Iter*` | the plain consumer |
| `Find` / `FindLast` | `Iter*` | `*T`, or `*Entry[K, V]` on a pair stream. `Find` short-circuits |
| `First` / `Last` | `Iter*` | `Last` drains |
| `Index` / `IndexValue(i, v)` | `Iter1*` | position in the stream, `-1` when absent |
| `Contains` / `ContainsValue(i, v)` | `Iter*` | short-circuits |
| `Every` | `Iter*` | true on an empty stream |
| `Count` / `CountBy` | `Iter*` | drains; a collection's `Len()` is free |
| `IsEmpty` | `Iter*`, every collection | pulls one element at most |
| `EqualBy(other, fn)` | `Iter*` | positional, so only meaningful on ordered sources |
| `Has` / `Len` / `At` / `IndexOf` | collections only | O(1), which is why they are not on iterators |

### Aggregate

| Operation | On | Notes |
| --------- | -- | ----- |
| `Reduce` / `ReduceX` | `Iter*` | seeded with the zero value |
| `Fold` / `FoldX` | `Iter*` | seeded with `init`; an error discards the accumulator |
| `MinBy` / `MaxBy` | `Iter*` | `*T`, nil on empty; a tie goes to the first element |
| `Min(i)` / `Max(i)` | funcs over `Iter1*` | need `T cmp.Ordered` |
| `SumBy` | `Iter*` | needs the exported `Number` constraint |
| `Sum(i)` | func over `Iter1*` | |

### Parallel

Each of these materializes its input, so they live on `Slice` and take an iterator
only as a convenience. `workers <= 0` means `runtime.NumCPU()`.

| Operation | On | Notes |
| --------- | -- | ----- |
| `ParallelMap` / `ParallelMapX` / `ParallelMapXC` | `Slice`, `Iter1*` | order preserved |
| `ParallelFilter` / `...X` / `...XC` | `Slice`, `Iter1*` | order preserved |
| `ParallelEach` / `...X` / `...XC` | `Slice`, `Iter1*` | side effects only |

The `X` forms return the first error and discard the partial result; it also cancels
work not yet started. The `XC` forms take a `context.Context` and pass each call a
context cancelled by that first error. The callback runs on many goroutines at once
and nothing here synchronizes what it touches.

## Technical decisions

### Generic methods, and their three limits

Verified against go1.27.0.

| | |
| --- | --- |
| A method may introduce new type parameters | yes |
| ...with any constraint (`comparable`, `cmp.Ordered`, ...) | yes |
| ...more than one | yes |
| ...on a pointer receiver of a generic struct | yes |
| Inference works when chaining them | yes |
| A method may **tighten the receiver's own** type parameter | **no** |
| A generic method may appear in an interface | **no** |
| A method may return its own generic type wrapped around its own parameter | **no** |

The last three explain every package-level function in the package.

**Tightening.** `Slice[T any]` cannot grow a `.Set()`, because that needs
`T comparable` and `T` is already `any`. So `CollectSet`, `Dedup`, `Min`, `Max`,
`Sum`, `ContainsValue` and `IndexValue` are functions. Each has a fluent twin that
takes a key function, moving the constraint onto a genuinely new parameter:

```go
CollectSet(s.Values())                       // function
s.Values().SetBy(func(t T) T { return t })   // generic method, K inferred as T
```

**No generic methods in interfaces.** This rules out a `Collection[T]` abstraction
over the five collections. They share conventions, not an interface.

**Instantiation cycles.** A method whose return type re-instantiates the receiver's
own generic type around a type built from its parameters makes the compiler chase
instantiations forever:

```text
Chunk as a method on Iter[T], returning Iter[Slice[T]]
  instantiating Iter[T] pulls in Iter[Slice[T]]
  which pulls in Iter[Slice[Slice[T]]]
  which pulls in ...
```

`go build` rejects that as `instantiation cycle`. It hit three times, in two shapes:

- **No method form exists.** `Chunk` and `Window` return `Iter[Slice[T]]` from
  `Iter[T]`. `GroupBy` returns `Map[K, Slice[T]]`, which cycles through `Map.Values()`
  handing back `Iter[Slice[T]]`. All functions.
- **A choice of which side to give up.** `Iter2.Pairs() Iter[Entry[K, V]]` and
  `Iter.KeyBy[K] Iter2[K, T]` each compile alone and cycle together, because the pair
  layer and the single-value layer then reference each other without end. One had to
  become a function; `KeyBy` earns the method, so `Pairs` is the function.

A free function is only instantiated where it is called, so it never forms the cycle.
The rule to carry forward: **an operation that wraps its elements in another generic
type cannot be a method.**

### One helper carries all the error plumbing

A fallible callback inside the `yield func(...) bool` protocol has nowhere to put its
error. Written out, every fallible operation repeats the same out-of-band capture and
merge — measured at 22 lines, about 45 times over. Hoisted into one unexported method
per fallible type:

```go
// each runs fn over the stream and merges a callback error with a stream abort.
func (i IterX[T]) each(fn func(T) (bool, error)) error {
	var ferr error
	err := i(func(t T) bool {
		ok, e := fn(t)
		if e != nil {
			ferr = e
			return false
		}
		return ok
	})
	if ferr != nil {
		return ferr
	}
	return err
}
```

Every fallible operation is then its own logic and nothing else:

```go
func (i IterX[T]) FilterX(fn func(T) (bool, error)) IterX[T] {
	return func(yield func(T) bool) error {
		return i.each(func(t T) (bool, error) {
			ok, err := fn(t)
			if err != nil {
				return false, err
			}
			if !ok {
				return true, nil
			}
			return yield(t), nil
		})
	}
}
```

Two `each` methods, 35 lines together, against roughly a thousand lines of repeated
plumbing.

### Where reuse pays, and where it does not

The rule is to implement the most complex version and reuse it, without chaining
reuses. Both halves hold, for different operations.

**Terminal operations** follow it as written. `Slice`, `Collect*`, `Find`, `Reduce`,
`Count`: the fallible version is the real implementation and the infallible one
discards a nil error, exactly as `Iter.Slice()` delegates to `IterX.Slice()`.

**Lazy operations deliberately do not.** `Iter.Filter` could delegate:

```go
return i.X().FilterX(fn).Iter(nil)   // one line, three closure layers per element
```

or be written directly:

```go
return func(yield func(T) bool) {    // eight lines, no adapters in the hot path
	for t := range i {
		if fn(t) && !yield(t) {
			return
		}
	}
}
```

The direct form costs seven lines and saves two wrappers on every element of every
stream, in the place where a stack trace is already hard to read. So `Iter` and
`Iter2` get direct implementations, `IterX` and `Iter2X` build on `each`, and
delegation depth never exceeds one.

This is argued from call-stack depth, not from measurement. See
[Known gaps](#known-gaps).

### Which callbacks get an X form

A callback gets one when **its result reaches the output**: `Map`, `MapKeys`,
`MapValues`, `KeyBy`, `Collapse`, `Scan`, `Fold`, `FlatMap`, and every predicate that
decides what survives — `Filter`, `TakeWhile`, `DropWhile`, `Find`, `Contains`,
`Every`, `Index`, `CountBy`.

It does not when it is a **key extractor used only internally**: `DedupBy`,
`CompactBy`, `SortBy`, `GroupBy`, `MinBy`, `MaxBy`, `SumBy`, `EqualBy`. A key derived
from an element for comparing, grouping or deduplicating it rarely fails, and an `X`
form for each would double the family for very little.

Where such a key genuinely can fail, the composition is already there:

```go
i.KeyByX(parseKey).MinBy(func(k int, _ T) int { return k })   // fallible key
SumX(i.MapX(strconv.Atoi))                                    // fallible value
```

### State lives inside the returned closure

Every lazy operation that carries state — `Take`, `Drop`, `DropWhile`, `Enumerate`,
`Scan`, `DedupBy`, `Chunk`, `Window`, `Zip` — declares it **inside** the function it
returns, never beside it. An iterator is a function and may be called more than once;
state one line too far up is a silent bug on the second iteration. This is enforced by
a test, and three deliberate mutations were used to confirm the test catches it.

### Two aborts have no single winner

`Zip` and `EqualBy` take an infallible other side, even on the `X` types. Zipping two
fallible streams would leave two aborts to reconcile with no principled answer for
which one wins, so the API does not offer the choice.

### What is deliberately absent

| Not there | Why |
| --------- | --- |
| `Reject` | `!fn(t)` at the call site beats a second name for `Filter` |
| `Some` | it would be a second name for `Contains` |
| `All(fn)` as a predicate | `All()` already means the pair iterator, so the predicate is `Every` |
| `Clear` on `Slice` | a slice is a value; `Slice[T]{}` at the call site says it better |
| `Grow` on `Map` / `Set` | Go offers no way to add capacity to a map once made; size with `make` |
| `Merge` / `InvertBy` on `OrderedMap` | merging into an ordered map raises a question — do incoming keys append, or keep the other map's order — that should be answered deliberately, not silently |
| an `I` (index) variant of anything | `Enumerate()` turns any stream into a pair stream, which covers it once instead of forty times |
| a `Collection[T]` interface | generic methods cannot appear in interfaces |

### Positional access hands back copies, except on Slice

`Slice.At(n)` returns a pointer **into** the slice, so writing through it changes the
element. `OrderedMap.At` and `OrderedSet.At` return a **copy**: those types keep an
index keyed on the value, and writing a new key through a pointer into the order would
leave the index describing something that is no longer there.

`Find`, `MinBy` and `MaxBy` also return a pointer to a copy — an iterator has no
addressable source to point into.

### Delete on the ordered types is O(n)

`OrderedMap.Delete` and `OrderedSet.Remove` repair every later index, so deleting in a
loop over a large collection is quadratic. Rebuild with a `Filter` and
`CollectOrdered*` instead.

### Package layout: by operation, not by type

The types live in `type_*.go`, one per type, holding the declaration, its conversions
and its entry points — nothing else. Every operation lives in `op_*.go` grouped by
what it does, so all eight `Filter` and `FilterX` methods across all four iterator
types sit in one 81-line file and the delegation between them is visible at once.
Collection-only methods are in `struct_*.go`.

The prefixes exist because `Map` is both a type and an operation, so `map.go` could
not be both. A method must also be declared in the same package as its receiver, so
there are no subpackages: the directory is flat and wide, the way `slices` and
`strings` are.

### A new module path instead of a /v2 suffix

This API replaced an earlier one, `go-slicest`, whose name had stopped describing what
the package was. Renaming the module rather than versioning it turns out to be the
cleaner answer to the compatibility question as well:

- `github.com/jannes-dailidow/go-slicest` keeps serving `v0.1.x` to anyone already on
  it, unchanged and forever. Nothing breaks, because nothing moved.
- `github.com/jannes-dailidow/go-collections` starts fresh. There is no earlier
  version at this path, so no `/vN` suffix is needed and no compatibility promise is
  inherited.

Had this API been published at the *same* path as the old one, it would have relied on
the module still being `v0`, where semver permits breaking changes and Go requires no
path suffix. Past `v1.0.0` that route closes: it would have had to be
`.../go-slicest/v2`, in a `v2/` subdirectory with its own `go.mod`.

The package is `collections` while the last path element is `go-collections`. That
mismatch is conventional in Go — the `go-` prefix disambiguates the repository, not the
package — and the compiler reads the package clause, not the path.

## Migrating from v0.1

The import path changed, so nothing breaks until you choose to move:

```go
import "github.com/jannes-dailidow/go-slicest"      // v0.1, still works
import "github.com/jannes-dailidow/go-collections"  // this
```

The old API was a set of package-level functions over `~[]T`, with suffixes combined
in order: `X` for an error, `I` for an index, `C` for a context, `Value` for a plain
value instead of a callback. Most of that dissolved into the type system.

| Old suffix | Now |
| ---------- | --- |
| `X` | the `IterX` / `Iter2X` types, or an `X` method for a fallible callback |
| `I` | `Enumerate()`, which turns any stream into `Iter2[int, T]` |
| `C` | kept, on the parallel operations only |
| `Value` | kept, as `ContainsValue` and `IndexValue` |

So `MapXI(s, fn)` is now `s.Values().Enumerate().MapX(fn)`, and the four functions
`Map`, `MapX`, `MapI`, `MapXI` are one method per iterator type.

| Old | New |
| --- | --- |
| `Map(s, fn)` | `s.Values().Map(fn)` |
| `Filter(s, fn)` | `s.Values().Filter(fn)` |
| `Find(s, fn)` | `s.Values().Find(fn)` — returns a pointer to a **copy**, not into `s` |
| `Contains(s, fn)` / `ContainsValue` | `s.Values().Contains(fn)` / `ContainsValue(s.Values(), v)` |
| `Index(s, fn)` / `IndexValue` | `s.Values().Index(fn)` / `IndexValue(s.Values(), v)` |
| `Reduce(s, fn)` | `s.Values().Reduce(fn)` — **argument order changed**, see below |
| `ReduceD(s, init, fn)` | `s.Values().Fold(init, fn)` |
| `Flatten(ss)` | `FlattenSlices(ss.Values())` |
| `ToMap(s, fn)` | `s.MapBy(fn)` |
| `FromMap(m, fn)` | `m.All().Collapse(fn)` |
| `MapKeys(m)` / `MapValues(m)` | `m.Keys()` / `m.Values()` |
| `Deduplicate(s)` / `DeduplicateFunc(s, fn)` | `Dedup(s.Values())` / `s.Values().DedupBy(fn)` |
| `Intersection(a, b)` | `CollectSet(a.Values()).Intersect(CollectSet(b.Values()))` |
| `Unique(a, b)` | the same, with `SymDiff` — the old name meant symmetric difference |
| `ParallelMap(s, fn, n)` | `s.ParallelMap(fn, n)` |

**One breaking change that the compiler will not catch.** `Reduce` and `Fold` take the
accumulator **first** — `fn(acc, t)` — matching `Scan`, the stdlib and every other
language's fold. The old `Reduce` took the element first, `fn(t, acc)`. When both are
the same type, swapping them compiles and gives the wrong answer.

## Testing

Every operation has its own tests, and two shared conformance suites carry the
invariants that have to hold everywhere.

**`conformance_test.go`** registers each lazy operation configured as a pass-through —
`Filter` accepting everything, `Map` applying identity — so the expected output is
always the source. Around thirty table lines produce seventy-five subtests. Adding an
operation should cost a line, not another six tests.

| Invariant | What it catches |
| --------- | --------------- |
| Passes every element through, in order | the obvious case |
| Stays lazy | an accidental `Slice()` inside a lazy operation |
| A consumer `break` leaves the error nil | conflating a break with an abort while merging errors |
| A callback error aborts the stream | the `each` contract |
| An upstream abort propagates | an operation that ignores the error `each` returns |
| An empty source yields nothing | |
| It can be iterated twice | state kept beside the returned closure instead of inside it |

**`conformance_terminal_test.go`** registers thirty-five terminal operations against
two more: an abort on the very first pull surfaces from every one of them, and a clean
stream never invents an error.

Both suites were checked by mutation — deliberately breaking one operation and
confirming the table fails.

Three things the suites cannot cover, tested individually instead:

- **Buffering operations.** `SortBy`, `SortFunc` and `Reverse` must see everything
  before yielding, so registering them as pass-throughs would hang the laziness test
  on its endless source.
- **Truncating operations.** A pass-through cannot expose leaked state in a `Take`
  whose `n` is at least the stream length, because it never reaches its limit. Those
  live in `op_take_test.go` with truncating configurations.
- **Shape-changing operations.** `KeyBy`, `Collapse`, `Enumerate`, `Chunk`, `Window`,
  `GroupBy` and `Partition` change the type of the stream, so no single table holds
  them.

## Project layout

```text
doc.go                  package documentation
constraints.go          Number
helpers.go              each, the shared error plumbing

type_iter.go            each type: declaration, conversions, entry points
type_iter2.go
type_iter_x.go
type_iter2_x.go
type_slice.go
type_map.go
type_set.go
type_ordered_map.go
type_ordered_set.go
type_entry.go

op_map.go               Map, MapX, MapKeys, MapValues  -- all four iterator types
op_filter.go            Filter, FilterX
op_take.go              Take, Drop, TakeWhile, DropWhile
op_dedup.go             DedupBy, CompactBy, Dedup, Compact
op_pair.go              KeyBy, Collapse, Enumerate, Zip, Pairs
op_flatten.go           FlatMap, Flatten, FlattenSlices
op_chunk.go             Chunk, Window
op_scan.go              Scan
op_each.go              Each, EachX, Tap
op_sort.go              SortBy, SortFunc, Reverse
op_group.go             GroupBy, GroupByOrdered, Partition
op_find.go              Find, FindLast, First, Last
op_index.go             Index, IndexValue
op_contains.go          Contains, ContainsValue, Every
op_count.go             Count, CountBy, IsEmpty
op_reduce.go            Reduce, Fold
op_minmax.go            MinBy, MaxBy, Min, Max
op_sum.go               SumBy, Sum
op_equal.go             EqualBy
op_collect.go           SetBy, MapBy

struct_slice.go         collection-only methods, per collection
struct_map.go
struct_set.go
struct_ordered_map.go
struct_ordered_set.go
struct_setops.go        Union, Intersect, Diff, SymDiff, IsSubset  -- both set types

parallel_map.go         ParallelMap and its X and C variants
parallel_filter.go
parallel_each.go
```

Conventions carried over from the previous version and still followed: a variable
naming a generic uses the same letter in lowercase, a variable collecting a result is
called `result`, slice results are preallocated when the final size is known and never
guessed, and functions run simple to complex within a file.

## Known gaps

- **The parallel family has never run under the race detector.** `-race` needs cgo and
  a C compiler. The tests pass repeatedly, including thirty consecutive runs of the
  parallel file, and the concurrency is a direct port of the previous version's — but
  repetition is not the race detector. Run `CGO_ENABLED=1 go test -race ./...` on a
  machine with gcc.
- **Nothing is benchmarked.** The direct-versus-delegated decision for lazy operations
  was argued from call-stack depth, not measured.
- **The conventions here are not all in godoc.** `doc.go` carries the naming rules and
  the implementation rule; the rest of this document is not reachable from
  `go doc`.
