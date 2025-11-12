package main

import (
	"fmt"
)

type traveler struct {
	name   string
	age    int
	vision string
}

type whale struct {
	name  string
	age   int
	hobby string
}

type tourist struct {
	name string
}

type travelerconf struct {
	name   string
	age    int
	vision string
}

type whaleconf struct {
	name  string
	age   int
	hobby string
}

type touristconf struct {
	name string
}

type player interface {
	traveler | whale | tourist
}

type playerconf interface {
	*travelerconf | *whaleconf | *touristconf
}

func (t *traveler) init(cfg *travelerconf) {
	t.name = cfg.name
	t.age = cfg.age
	t.vision = cfg.vision
}

func (w *whale) init(cfg *whaleconf) {
	w.name = cfg.name
	w.age = cfg.age
	w.hobby = cfg.hobby
}

func (t *tourist) init(cfg *touristconf) {
	t.name = cfg.name
}

func create[T player, C playerconf](cfg C) T {
	var t T
	switch v := any(&t).(type) {
	case *traveler:
		v.init(any(cfg).(*travelerconf))
	case *whale:
		v.init(any(cfg).(*whaleconf))
	case *tourist:
		v.init(any(cfg).(*touristconf))
	}
	return t
}

func main() {
	a := create[traveler](&travelerconf{
		name:   "Alice",
		age:    30,
		vision: "Explore the world",
	})
	b := create[whale](&whaleconf{
		name:  "Willy",
		age:   5,
		hobby: "Jumping",
	})
	c := create[tourist](&touristconf{
		name: "Bob",
	})
	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(c)
}
