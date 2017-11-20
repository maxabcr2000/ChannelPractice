package character

import (
	"fmt"
	"time"
)

type SourceStage struct {
	OutboundCh chan Character
}

func NewSourceStage(outCh chan Character) *SourceStage {
	return &SourceStage{
		OutboundCh: outCh,
	}
}

func (stage *SourceStage) Start() {
	i := 0
	for {
		name := fmt.Sprintf("TestPlayer%d", i)
		stage.OutboundCh <- NewCharacter(i, name)
		fmt.Printf("New player: %s\n", name)
		i++
		time.Sleep(1 * time.Second)
	}
}
