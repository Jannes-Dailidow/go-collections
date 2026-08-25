package collections

import (
	"errors"
	"strconv"
	"testing"
)

func TestMapChangesTheElementType(t *testing.T) {
	got := ints(3).Map(strconv.Itoa).Native()
	if len(got) != 3 || got[0] != "0" || got[2] != "2" {
		t.Fatalf("got %v", got)
	}
}

func TestMapXReturnsTheCallbackError(t *testing.T) {
	got, err := ints(6).MapX(func(v int) (string, error) {
		if v == 2 {
			return "", errBoom
		}
		return strconv.Itoa(v), nil
	}).Native()
	if !errors.Is(err, errBoom) {
		t.Fatalf("want errBoom, got %v", err)
	}
	if got != nil {
		t.Fatalf("the partial result must be discarded, got %v", got)
	}
}

func TestMapXParsesIntoTheFallibleFamily(t *testing.T) {
	src := Slice[string]{"1", "2", "3"}
	got, err := src.Values().MapX(strconv.Atoi).Native()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertInts(t, got, 1, 2, 3)
}

func TestMapXSurfacesAParseFailure(t *testing.T) {
	src := Slice[string]{"1", "nope", "3"}
	if _, err := src.Values().MapX(strconv.Atoi).Native(); err == nil {
		t.Fatal("want a parse error, got nil")
	}
}

func TestPairMapReplacesBothHalves(t *testing.T) {
	got := CollectMap(pairs(3).Map(func(k, v int) (string, bool) {
		return strconv.Itoa(k), v > 0
	}))
	if len(got) != 3 || got["0"] != false || got["2"] != true {
		t.Fatalf("got %v", got)
	}
}

func TestMapKeysChangesTheKeyType(t *testing.T) {
	got := CollectMap(pairs(3).MapKeys(func(k, _ int) string {
		return strconv.Itoa(k)
	}))
	if len(got) != 3 || got["1"] != 10 {
		t.Fatalf("got %v", got)
	}
}

// The key function sees the value too, which is the whole reason it takes both.
func TestMapKeysCanDeriveTheKeyFromTheValue(t *testing.T) {
	got := CollectMap(pairs(3).MapKeys(func(_, v int) int { return v }))
	if len(got) != 3 || got[20] != 20 {
		t.Fatalf("got %v", got)
	}
}

func TestMapValuesChangesTheValueType(t *testing.T) {
	got := CollectMap(pairs(3).MapValues(func(_, v int) string {
		return strconv.Itoa(v)
	}))
	if len(got) != 3 || got[2] != "20" {
		t.Fatalf("got %v", got)
	}
}

func TestMapValuesXReturnsTheCallbackError(t *testing.T) {
	_, err := CollectMapX(pairs(4).MapValuesX(func(k, v int) (int, error) {
		if k == 2 {
			return 0, errBoom
		}
		return v, nil
	}))
	if !errors.Is(err, errBoom) {
		t.Fatalf("want errBoom, got %v", err)
	}
}

func TestMapOnAMapCollection(t *testing.T) {
	m := Map[string, int]{"a": 1, "b": 2}
	got := CollectMap(m.All().MapValues(func(_ string, v int) int { return v * 10 }))
	if got["a"] != 10 || got["b"] != 20 {
		t.Fatalf("got %v", got)
	}
}

// A chain across both layers: values out of a map, mapped, collected in order.
func TestMapChainsThroughTheLayers(t *testing.T) {
	src := Slice[int]{3, 1, 2, 1}
	got := CollectOrderedSet(src.Values().
		Filter(func(v int) bool { return v > 1 }).
		Map(func(v int) int { return v * 100 }))
	assertInts(t, got.Native(), 300, 200)
}
