package character_test

import (
	"testing"

	"github.com/maxabcr2000/ChannelPractice/character"
)

var (
	completeCh = make(chan character.Action)
	failCh     = make(chan character.Action)
)

func InitCharManagementStage() *character.CharManagementStage {
	stage := character.NewCharManagementStage(completeCh, failCh)
	go stage.Start()
	return stage
}

func TestRegister(t *testing.T) {
	stage := InitCharManagementStage()
	action := character.NewAction("TestAction", "TestAction")
	go stage.Register(action)

	select {
	case <-completeCh:
		t.Error("Should return error when non-struct data was passed in.")
	case <-failCh:
	}

	inChar := character.NewCharacter(1, "TestBob")
	action.SetData(inChar)

	go stage.Register(action)

	select {
	case <-completeCh:
		outChar, ok := stage.TryGetCharacter(1)
		if !ok || inChar.Name != outChar.Name || inChar.ID != outChar.ID {
			t.Error("Register failed.")
		}
	case <-failCh:
	}

	go stage.Register(action)

	select {
	case <-completeCh:
		t.Error("Should return error when trying to register character with same id.")
	case <-failCh:
	}
}

func TestDelete(t *testing.T) {

}

func TestUpdate(t *testing.T) {

}

func TestRead(t *testing.T) {

}

func TestTryGetCharacter(t *testing.T) {

}
