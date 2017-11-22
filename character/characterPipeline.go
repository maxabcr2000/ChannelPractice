package character

import "strconv"
import "errors"
import "time"

const (
	DEFAULT_TIMEOUT = time.Duration(500 * time.Millisecond)
)

var (
	currentActionID = 0
)

type Action struct {
	ID          int
	Description string
	data        interface{}
	returnCh    chan Action
}

func NewAction(desc string, data interface{}) Action {
	currentActionID++
	return Action{
		ID:          currentActionID,
		Description: desc,
		data:        data,
	}
}

func (a *Action) GetDataAsString() (string, error) {
	stringVal, ok := a.data.(string)
	if !ok {
		return "", errors.New("Occurred error when converting data to string type")
	}

	return stringVal, nil
}

func (a *Action) GetDataAsInt() (int, error) {
	intVal, ok := a.data.(int)
	if !ok {
		return 0, errors.New("Occurred error when converting data to string type")
	}

	return intVal, nil
}

func (a *Action) GetDataAsCharacter() (Character, error) {
	structVal, ok := a.data.(Character)
	if !ok {
		return Character{}, errors.New("Occurred error when converting data to struct type: Character")
	}

	return structVal, nil
}

func (a *Action) GetDataAsError() (error, error) {
	errVal, ok := a.data.(error)
	if !ok {
		return nil, errors.New("Occurred error when converting data to error type")
	}

	return errVal, nil
}

func (a *Action) SetData(data interface{}) {
	a.data = data
}

type CharacterPipeline struct {
	completeCh       chan Action
	failCh           chan Action
	ErrorStage       *ErrorStage
	SinkStage        *SinkStage
	charManagerStage *CharManagementStage
	mockSourceStage  *MockSourceStage
}

func NewCharPipeline() *CharacterPipeline {
	return &CharacterPipeline{}
}

func (p *CharacterPipeline) Start() {
	p.completeCh = make(chan Action)
	p.failCh = make(chan Action)

	p.ErrorStage = NewErrorStage(p.failCh)
	p.SinkStage = NewSinkStage(p.completeCh)
	p.charManagerStage = NewCharManagementStage(p.completeCh, p.failCh)
	p.mockSourceStage = NewMockSourceStage(p.charManagerStage.GetRegisterCh())

	//go p.mockSourceStage.Start()
	go p.charManagerStage.Start()
	go p.SinkStage.Start()
	go p.ErrorStage.Start()
}

func (p *CharacterPipeline) RegisterWithTimeout(action Action, timeout time.Duration) bool {
	return p.register(action, timeout)
}

func (p *CharacterPipeline) Register(action Action) bool {
	return p.register(action, DEFAULT_TIMEOUT)
}

func (p *CharacterPipeline) register(action Action, timeout time.Duration) bool {
	for {
		select {
		case p.charManagerStage.GetRegisterCh() <- action:
			return false
		case <-time.After(timeout):
			return true
		}
	}
}

func (p *CharacterPipeline) DeleteWithTimeout(action Action, timeout time.Duration) bool {
	return p.delete(action, timeout)
}

func (p *CharacterPipeline) Delete(action Action) bool {
	return p.delete(action, DEFAULT_TIMEOUT)
}

func (p *CharacterPipeline) delete(action Action, timeout time.Duration) bool {
	stringVal, err := action.GetDataAsString()
	if err != nil {
		action.Description = err.Error()
		action.data = err
		p.failCh <- action
		return false
	}

	var intVal int
	intVal, err = strconv.Atoi(stringVal)
	if err != nil {
		action.Description = err.Error()
		action.data = err
		p.failCh <- action
		return false
	}

	action.data = intVal

	for {
		select {
		case p.charManagerStage.GetDeleteCh() <- action:
			return false
		case <-time.After(timeout):
			return true
		}
	}

}

func (p *CharacterPipeline) UpdateWithTimeout(action Action, timeout time.Duration) bool {
	return p.update(action, timeout)
}

func (p *CharacterPipeline) Update(action Action) bool {
	return p.update(action, DEFAULT_TIMEOUT)
}

func (p *CharacterPipeline) update(action Action, timeout time.Duration) bool {
	for {
		select {
		case p.charManagerStage.GetUpdateCh() <- action:
			return false
		case <-time.After(timeout):
			return true
		}
	}
}

func (p *CharacterPipeline) ReadWithTimeout(action Action, timeout time.Duration) bool {
	return p.read(action, timeout)
}

func (p *CharacterPipeline) Read(action Action) bool {
	return p.read(action, DEFAULT_TIMEOUT)
}

func (p *CharacterPipeline) read(action Action, timeout time.Duration) bool {
	stringVal, err := action.GetDataAsString()
	if err != nil {
		action.Description = err.Error()
		action.data = err
		p.failCh <- action
		return false
	}

	var intVal int
	intVal, err = strconv.Atoi(stringVal)
	if err != nil {
		action.Description = err.Error()
		action.data = err
		p.failCh <- action
		return false
	}

	action.data = intVal
	for {
		select {
		case p.charManagerStage.GetQueryCh() <- action:
			return false
		case <-time.After(timeout):
			return true
		}
	}

}

// func (p *CharacterPipeline) handleFailedAction(action Action, err error) {
// 	if err != nil {
// 		action.Description = err.Error()
// 		action.data = err
// 		p.failCh <- action
// 		return
// 	}
// }
