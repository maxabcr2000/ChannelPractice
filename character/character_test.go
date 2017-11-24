package character_test

import (
	"testing"

	"github.com/maxabcr2000/ChannelPractice/character"
)

func TestNewCharacter(t *testing.T) {
	char := character.NewCharacter(1, "TestBob")
	if char.Name != "TestBob" ||
		char.ID != 1 ||
		char.State != character.Normal ||
		char.Hp != 100 ||
		char.Mp != 100 ||
		char.Level != 1 ||
		char.AttackPower != 10 {
		t.Error("Create new character with error!")
	}
}
