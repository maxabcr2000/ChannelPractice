package character

import "fmt"

type ErrorStage struct {
	failedActions map[int]Action
	inboundCh     chan Action
}

func NewErrorStage(inCh chan Action) *ErrorStage {
	return &ErrorStage{
		failedActions: make(map[int]Action),
		inboundCh:     inCh,
	}
}

func (stage *ErrorStage) Start() {
	for {
		select {
		case action := <-stage.inboundCh:
			fmt.Println(action)
			stage.failedActions[action.ID] = action
		}
	}
}

func (stage *ErrorStage) CheckFailedAction(actionId int) bool {
	_, ok := stage.failedActions[actionId]
	return ok
}
