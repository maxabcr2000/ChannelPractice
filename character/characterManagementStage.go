package character

import "fmt"

type CharManagementStage struct {
	charMap    map[int]Character
	registerCh chan Action
	deleteCh   chan Action
	updateCh   chan Action
	queryCh    chan Action
	completeCh chan Action
	failCh     chan Action
}

func NewCharManagementStage(completeCh chan Action, failCh chan Action) *CharManagementStage {
	return &CharManagementStage{
		charMap:    make(map[int]Character),
		registerCh: make(chan Action),
		deleteCh:   make(chan Action),
		updateCh:   make(chan Action),
		queryCh:    make(chan Action),
		completeCh: completeCh,
		failCh:     failCh,
	}
}

func (c *CharManagementStage) Register(action Action) {
	c.registerCh <- action
}

func (c *CharManagementStage) Delete(action Action) {
	c.deleteCh <- action
}

func (c *CharManagementStage) Update(action Action) {
	c.updateCh <- action
}

// func (c *CharManagementStage) Read(char Character) {
// 	c.queryCh <- char
// }

func (c *CharManagementStage) Read(id int) Character {
	charVal, ok := c.charMap[id]
	if ok {
		return charVal
	}
	return Character{ID: 0, Name: "Nil"}
}

func (c *CharManagementStage) GetRegisterCh() chan Action {
	return c.registerCh
}

func (c *CharManagementStage) GetDeleteCh() chan Action {
	return c.deleteCh
}

func (c *CharManagementStage) GetUpdateCh() chan Action {
	return c.updateCh
}

func (c *CharManagementStage) GetQueryCh() chan Action {
	return c.queryCh
}

func (stage *CharManagementStage) Start() {
	// fmt.Println("Sleep for 10 seconds.")
	// time.Sleep(10 * time.Second)
	// fmt.Println("Sleep finished.")

	for {
		select {
		case action := <-stage.registerCh:
			char, err := action.GetDataAsCharacter()
			if err != nil {
				action.Description = err.Error()
				action.SetData(err)
				stage.failCh <- action
				continue
			}

			_, ok := stage.charMap[char.ID]

			if !ok {
				stage.charMap[char.ID] = char
				action.Description = "Complete character register"
				stage.completeCh <- action
			} else {
				err = fmt.Errorf("Character with id: %d already existed", char.ID)
				action.Description = err.Error()
				action.SetData(err)
				stage.failCh <- action
			}
		case action := <-stage.deleteCh:
			charID, err := action.GetDataAsInt()
			if err != nil {
				action.Description = err.Error()
				action.SetData(err)
				stage.failCh <- action
				continue
			}

			delete(stage.charMap, charID)
			action.Description = "Complete character delete"
			stage.completeCh <- action
		case action := <-stage.updateCh:
			char, err := action.GetDataAsCharacter()
			if err != nil {
				action.Description = err.Error()
				action.SetData(err)
				stage.failCh <- action
				continue
			}

			stage.charMap[char.ID] = char
			action.Description = "Complete character update"
			stage.completeCh <- action
		case action := <-stage.queryCh:
			charID, err := action.GetDataAsInt()
			if err != nil {
				action.Description = err.Error()
				action.SetData(err)
				stage.failCh <- action
				continue
			}

			charVal, ok := stage.charMap[charID]

			if ok {
				action.SetData(charVal)
				action.Description = "Complete character query"
				stage.completeCh <- action
			} else {
				err = fmt.Errorf("Request character with id: %d is not existed", charID)
				action.Description = err.Error()
				action.SetData(err)
				stage.failCh <- action
			}
		}
	}
}
