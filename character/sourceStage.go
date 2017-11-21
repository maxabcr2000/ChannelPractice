package character

import (
	"fmt"
	"time"
)

type SourceStage struct {
	outboundCh chan Character
}

func NewSourceStage(outCh chan Character) *SourceStage {
	return &SourceStage{
		outboundCh: outCh,
	}
}

func (stage *SourceStage) Start() {
	i := 1
	for {
		name := fmt.Sprintf("TestPlayer%d", i)
		stage.outboundCh <- NewCharacter(i, name)
		fmt.Printf("New player: %s\n", name)
		i++
		time.Sleep(1 * time.Second)
	}
}
