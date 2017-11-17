package character

import (
	"time"
)

const (
	Normal StateType = iota
	Defeated
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

	HorizontalMovement chan int
	VerticalMovement   chan int
	Damaged            chan int
	GainExp            chan int
	StateChange        chan StateType
}

func (char *character) EventSystem() {
	if char == nil {
		return
	}

	char.HorizontalMovement = make(chan int)
	char.VerticalMovement = make(chan int)
	char.Damaged = make(chan int)
	char.GainExp = make(chan int)
	char.StateChange = make(chan StateType)

	for {
		select {
		case moveX := <-char.HorizontalMovement:
			char.PosX += moveX
		case moveY := <-char.VerticalMovement:
			char.PosY += moveY
		case damage := <-char.Damaged:
			char.Hp -= damage
			if char.Hp <= 0 {
				char.TriggerStateChange(Defeated)
			}
		case exp := <-char.GainExp:
			char.Exp += exp
		case state := <-char.StateChange:
			char.State = state
		}

		time.Sleep(1 * time.Second)
	}
}

func (char *character) TriggerStateChange(state StateType) {
	if state > Poisoned {
		return
	}
	char.StateChange <- state
}

func NewCharacter(name string) *character {
	c := &character{Name: name, Hp: 100, Mp: 100, State: Normal, Level: 1, AttackPower: 10}
	go c.EventSystem()
	return c
}
