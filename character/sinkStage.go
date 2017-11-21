package character

import "fmt"

type SinkStage struct {
	inboundCh chan Character
}

func NewSinkStage(inCh chan Character) *SinkStage {
	return &SinkStage{
		inboundCh: inCh,
	}
}

func (stage *SinkStage) Start() {
	for {
		select {
		case char := <-stage.inboundCh:
			fmt.Println(char)
		}
	}
}
