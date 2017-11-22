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
			fmt.Printf("ErrorStage got action: %+v\n", action)
			stage.failedActions[action.ID] = action
		}
	}
}

func (stage *ErrorStage) CheckFailedAction(actionId int) (Action, bool) {
	action, ok := stage.failedActions[actionId]
	fmt.Println("ErrorStage.CheckFailedAction():", ok)
	return action, ok
}
