package character

type CharManagementStage struct {
	charMap    map[int]Character
	registerCh chan Character
	deleteCh   chan Character
	updateCh   chan Character
	queryCh    chan Character
	outboundCh chan Character
}

func NewCharManagementStage(outCh chan Character) *CharManagementStage {
	return &CharManagementStage{
		charMap:    make(map[int]Character),
		registerCh: make(chan Character),
		deleteCh:   make(chan Character),
		updateCh:   make(chan Character),
		queryCh:    make(chan Character),
		outboundCh: outCh,
	}
}

func (c *CharManagementStage) Register(char Character) {
	c.registerCh <- char
}

func (stage *CharManagementStage) Start() {
	for {
		select {
		case char := <-stage.RegisterCh:
			_, ok := stage.charMap[char.ID]
			if !ok {
				stage.charMap[char.ID] = char
				stage.OutboundCh <- char
			}
		case char := <-stage.DeleteCh:
			delete(stage.charMap, char.ID)
		case char := <-stage.UpdateCh:
			stage.charMap[char.ID] = char
		case char := <-stage.QueryCh:
			charVal, ok := stage.charMap[char.ID]
			if ok {
				stage.OutboundCh <- charVal
			}
		}
	}
}
