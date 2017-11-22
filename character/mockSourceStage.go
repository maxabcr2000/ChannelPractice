package character

import (
	"fmt"
	"time"
)

type MockSourceStage struct {
	outboundCh chan Action
}

func NewMockSourceStage(outCh chan Action) *MockSourceStage {
	return &MockSourceStage{
		outboundCh: outCh,
	}
}

func (stage *MockSourceStage) Start() {
	i := 1
	for {
		name := fmt.Sprintf("TestPlayer%d", i)
		stage.outboundCh <- NewAction("Create mock users", NewCharacter(i, name))
		fmt.Printf("New player: %s\n", name)
		i++
		time.Sleep(1 * time.Second)
	}
}
