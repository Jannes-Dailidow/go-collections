package collections

import (
	"errors"
	"strings"
	"testing"
)

func TestSetByDedupsOnTheKey(t *testing.T) {
	src := Slice[string]{"Ada", "ada", "Bob"}
	got := src.Values().SetBy(strings.ToLower)
	if got.Len() != 2 || !got.Has("ada") || !got.Has("bob") {
		t.Fatalf("got %v", got)
	}
}

// The identity case, which is what CollectSet does as a function.
func TestSetByWithAnIdentityKeyMatchesCollectSet(t *testing.T) {
	src := Slice[int]{3, 1, 3}
	fluent := src.Values().SetBy(func(v int) int { return v })
	asFunc := CollectSet(src.Values())
	if fluent.Len() != asFunc.Len() || !fluent.IsSubset(asFunc) {
		t.Fatalf("%v vs %v", fluent, asFunc)
	}
}

func TestSliceSetByAvoidsTheIterator(t *testing.T) {
	src := Slice[int]{1, 2, 2}
	if got := src.SetBy(func(v int) int { return v }); got.Len() != 2 {
		t.Fatalf("got %v", got)
	}
}

func TestMapByLastKeyWins(t *testing.T) {
	src := Slice[string]{"a1", "b2", "a3"}
	got := src.Values().MapBy(func(s string) (string, string) {
		return s[:1], s[1:]
	})
	if got.Len() != 2 {
		t.Fatalf("want 2, got %v", got)
	}
	if v, _ := got.Get("a"); v != "3" {
		t.Fatalf("the last a must win, got %q", v)
	}
}

func TestSliceMapByIsTheRootPackagesToMap(t *testing.T) {
	type user struct {
		id   int
		name string
	}
	src := Slice[user]{{1, "ada"}, {2, "bob"}}
	got := src.MapBy(func(u user) (int, string) { return u.id, u.name })
	if got[1] != "ada" || got[2] != "bob" {
		t.Fatalf("got %v", got)
	}
}

func TestSetByAndMapByOnAFallibleStream(t *testing.T) {
	if _, err := abortAt(2).SetBy(func(v int) int { return v }); !errors.Is(err, errBoom) {
		t.Fatalf("want errBoom, got %v", err)
	}
	if _, err := abortAt(2).MapBy(func(v int) (int, int) { return v, v }); !errors.Is(err, errBoom) {
		t.Fatalf("want errBoom, got %v", err)
	}

	got, err := ints(3).X().SetBy(func(v int) int { return v })
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Len() != 3 {
		t.Fatalf("got %v", got)
	}
}

// SetBy after a Map is the pattern that makes the constraint dance worth it.
func TestSetByAfterATransform(t *testing.T) {
	got := ints(6).
		Filter(func(v int) bool { return v > 1 }).
		Map(func(v int) int { return v % 2 }).
		SetBy(func(v int) int { return v })
	if got.Len() != 2 {
		t.Fatalf("got %v", got)
	}
}
