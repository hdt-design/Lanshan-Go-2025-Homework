package main

import (
	"fmt"
	"strings"
)

type alice struct {
	name  string
	age   int
	lover string
}

type bob struct {
	name  string
	age   int
	intro string
}

type carol struct {
	name string
}

type aliceconf struct {
	name  string
	age   int
	lover string
}

type bobconf struct {
	name string
	age  int
	tags []string
}

type carolconf struct {
	name string
}

type defaultptr[T usedtype, C usedtypeconf] interface {
	*T
	init(C)
}

type usedtype interface {
	alice | bob | carol
}

type usedtypeconf interface {
	*aliceconf | *bobconf | *carolconf
}

func (a *alice) init(cfg *aliceconf) {
	a.name = cfg.name
	a.age = cfg.age
	a.lover = cfg.lover
}

func (b *bob) init(cfg *bobconf) {
	b.name = cfg.name
	b.age = cfg.age
	b.intro = fmt.Sprintf("I like to %v", strings.Join(cfg.tags, ", "))
}

func (c *carol) init(cfg *carolconf) {
	c.name = cfg.name
}

func Default[T usedtype, C usedtypeconf, ptr defaultptr[T, C]](cfg C) T {
	var t T
	ptr.init(&t, cfg)
	return t
}

func main() {
	fmt.Println(Default[alice, *aliceconf](&aliceconf{
		name:  "Alice",
		age:   20,
		lover: "Bob",
	}))
	fmt.Println(Default[bob, *bobconf](&bobconf{
		name: "Bob",
		age:  18,
		tags: []string{"sports", "coding"},
	}))
}
