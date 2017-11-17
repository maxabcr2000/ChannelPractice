package character_test

import (
	"testing"
	"time"

	"github.com/maxabcr2000/ChannelPractice/character"
)

func TestNewCharacter(t *testing.T) {
	char := character.NewCharacter("TestBob")
	if char.Name != "TestBob" || char.State != character.Normal || char.Hp != 100 || char.Mp != 100 || char.Level != 1 || char.AttackPower != 10 {
		t.Error("Create new character with error!")
	}
}

func TestTriggerStateChange(t *testing.T) {
	char := character.NewCharacter("TestBob")
	char.TriggerStateChange(character.Frozen)
	time.Sleep(5 * time.Second)

	if char.State != character.Frozen {
		t.Errorf("char.State should be %d but get %d", character.Frozen, char.State)
	}
}
