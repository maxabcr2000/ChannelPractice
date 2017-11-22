package character

import "fmt"

type SinkStage struct {
	completeActions map[int]Action
	inboundCh       chan Action
}

func NewSinkStage(inCh chan Action) *SinkStage {
	return &SinkStage{
		completeActions: make(map[int]Action),
		inboundCh:       inCh,
	}
}

func (stage *SinkStage) Start() {
	for {
		select {
		case action := <-stage.inboundCh:
			fmt.Printf("SinkStage got action: %+v\n", action)
			stage.completeActions[action.ID] = action
		}
	}
}

func (stage *SinkStage) CheckCompletedAction(actionId int) (Action, bool) {
	action, ok := stage.completeActions[actionId]
	fmt.Println("SinkStage.CheckCompletedAction():", ok)
	return action, ok
}
