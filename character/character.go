package character

const (
	Normal StateType = iota
	Defeated
	Frozen
	Poisoned
)

type StateType int

type Character struct {
	ID          int
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

func NewCharacter(id int, name string) Character {
	c := Character{
		ID:          id,
		Name:        name,
		Hp:          100,
		Mp:          100,
		State:       Normal,
		Level:       1,
		AttackPower: 10,
	}
	return c
}

func (c *Character) Move() {
	c.PosX += 10
	c.PosY += 10
}

func (c *Character) TryCollide(target *Character) bool {
	if c.PosX == target.PosX && c.PosY == target.PosY {
		return true
	}
	return false
}

func (c *Character) Battle() {
	c.Hp -= 10
}
