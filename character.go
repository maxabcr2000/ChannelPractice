package main

import (
	"time"
)

const (
	Normal StateType = iota
	Frozen
	Poisoned
)

type StateType int

type character struct {
	Name        string
	PosX        int
	PosY        int
	Hp          int
	Mp          int
	State       StateType
	Exp         int
	Level       int
	AttackPower int
}

func (char *Character) EventSystem() {
	horizontalMovement := make(chan int)
	verticalMovement := make(chan int)
	damaged := make(chan int)
	gainExp := make(chan int)

	for {
		time.Sleep(1 * time.Second)

		select {
		case moveX := <-horizontalMovement:

		case moveY := <-verticalMovement:

		case damage := <-damaged:

		case exp := <-gainExp:

		}
	}
}

func (char *character) NewCharater(name string) *character {
	c := &character{Name: name, Hp: 100, Mp: 100, State: Normal, Level: 1, AttackPower: 10}

	return c
}
