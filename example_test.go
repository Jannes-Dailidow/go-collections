package collections_test

import (
	"fmt"
	"strconv"

	"github.com/jannes-dailidow/go-collections"
)

type user struct {
	Name    string
	Surname string
	Team    string
	Active  bool
	ID      string
}

func people() collections.Slice[user] {
	return collections.Slice[user]{
		{Name: "Ada", Surname: "Lovelace", Team: "analytics", Active: true, ID: "7"},
		{Name: "Bob", Surname: "Brown", Team: "platform", Active: false, ID: "3"},
		{Name: "Cyd", Surname: "Charisse", Team: "analytics", Active: true, ID: "1"},
		{Name: "Dot", Surname: "Aitken", Team: "platform", Active: true, ID: "9"},
	}
}

// A chain stays lazy until the Slice() at the end.
func Example_chain() {
	names := people().Values().
		Filter(func(u user) bool { return u.Active }).
		SortBy(func(u user) string { return u.Surname }).
		Map(func(u user) string { return u.Name }).
		Take(10).
		Slice()

	fmt.Println(names.Native())
	// Output: [Dot Cyd Ada]
}

// A fallible transform moves into the X family and carries one error out.
func Example_fallible() {
	ids, err := people().Values().
		MapX(func(u user) (int, error) { return strconv.Atoi(u.ID) }).
		Slice()

	fmt.Println(ids.Native(), err)

	// A single failure aborts the stream and discards the partial result.
	broken := collections.Slice[string]{"1", "nope", "3"}
	parsed, err := broken.Values().MapX(strconv.Atoi).Slice()
	fmt.Println(parsed, err != nil)

	// Output:
	// [7 3 1 9] <nil>
	// [] true
}

// GroupByOrdered keeps the groups in the order their keys first appeared.
func Example_group() {
	byTeam := collections.GroupByOrdered(people().Values(), func(u user) string {
		return u.Team
	})

	for team, members := range byTeam.All() {
		fmt.Println(team, members.Len())
	}
	// Output:
	// analytics 2
	// platform 2
}

// The collections carry the O(1) questions an iterator cannot answer cheaply.
func Example_collections() {
	teams := collections.CollectOrderedSet(people().Values().Map(func(u user) string {
		return u.Team
	}))

	fmt.Println(teams.Native(), teams.Has("platform"), teams.IndexOf("platform"))
	// Output: [analytics platform] true 1
}
