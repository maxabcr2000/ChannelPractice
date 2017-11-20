package character

import "fmt"

type SinkStage struct {
	InboundCh chan Character
}

func NewSinkStage(inCh chan Character) *SinkStage {
	return &SinkStage{
		InboundCh: inCh,
	}
}

func (stage *SinkStage) Start() {
	for {
		select {
		case char := <-stage.InboundCh:
			fmt.Println(char)
		}
	}
}
