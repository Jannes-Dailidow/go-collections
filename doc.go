// Package collections provides generic collection types with a shared iterator layer,
// built on Go's iter package.
//
// # Layers
//
// Five collections -- [Slice], [Map], [Set], [OrderedMap] and [OrderedSet] -- and
// four iterator types -- [Iter], [Iter2], [IterX] and [Iter2X]. The collections
// carry only what they can do better than a generic iterator: O(1) lookup, set
// algebra and positional access. Every other operation lives on the iterator types,
// one call away through Values, All or Keys.
//
// # Naming
//
//	Native             the plain Go type: []T, map[K]V, map[T]struct{}
//	Seq, Seq2          the stdlib iter.Seq and iter.Seq2 equivalent
//	All, Keys, Values  this package's Iter and Iter2
//	Collect*           builds a collection by draining an iterator
//	*By                takes a key function, so a constraint can land on the key
//
// # The X suffix
//
// On a type it means the stream can abort: [IterX] is func(yield func(T) bool) error.
// On a method it means the callback can fail, which lifts an infallible stream into
// the fallible family, so [Iter.MapX] returns an [IterX].
//
// An abort is not a value in the stream. It stops iteration and surfaces once at the
// end, and collectors discard the partial result. A consumer break is not an abort
// and leaves the error nil.
//
// # Functions rather than methods
//
// Eleven operations here are functions rather than methods, for two reasons.
//
// A method cannot tighten its receiver's own type parameter: Slice[T any] cannot grow
// a Set() that needs T comparable. That accounts for Collect*, Dedup, Compact, Min,
// Max, Sum, ContainsValue and IndexValue. Where a fluent form is wanted, a *By method
// takes a key function so the constraint lands on the key instead, as in
// s.Values().SetBy(fn).
//
// A method also cannot return its own generic type wrapped around its own parameter,
// because the compiler then chases instantiations without end -- a Chunk method on
// Iter[T] returning Iter[Slice[T]] needs Iter[Slice[Slice[T]]], and so on. That
// accounts for [Chunk], [Window], [GroupBy], [Pairs] and [Flatten]. A free function
// is only instantiated where it is called, so it never forms the cycle. Anything
// added later that wraps its elements in another generic type has to be a function
// too.
//
// # Implementation
//
// Terminal operations -- Slice, Collect*, Find, Reduce -- are written once on the
// fallible type, and the infallible variant discards a nil error. [Iter.Slice] is
// the model.
//
// Lazy operations deliberately do not follow that rule. [Iter] and [Iter2] get
// direct implementations, because delegating through X() and Iter(nil) would add two
// closure layers to every element of every stream to save seven lines. [IterX] and
// [Iter2X] build on each, which merges a callback error with a stream abort in one
// place. Delegation depth never exceeds one.
package collections
