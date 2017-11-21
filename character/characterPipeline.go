package character

import "strconv"
import "fmt"

type CharacterPipeline struct {
	sinkStage        *SinkStage
	charManagerStage *CharManagementStage
	sourceStage      *SourceStage
}

func NewCharPipeline() *CharacterPipeline {
	return &CharacterPipeline{}
}

func (p *CharacterPipeline) Start() {
	outBoundCh := make(chan Character)

	p.sinkStage = NewSinkStage(outBoundCh)
	p.charManagerStage = NewCharManagementStage(outBoundCh)
	p.sourceStage = NewSourceStage(p.charManagerStage.GetRegisterCh())

	go p.sourceStage.Start()
	go p.charManagerStage.Start()
	go p.sinkStage.Start()
}

func (p *CharacterPipeline) Register(char Character) {
	p.charManagerStage.GetRegisterCh() <- char
}

func (p *CharacterPipeline) Delete(id string) {
	intVal, err := strconv.Atoi(id)
	if err != nil {
		fmt.Printf("Occurred error when converting character id.")
		return
	}
	char := NewCharacter(intVal, "")
	p.charManagerStage.GetDeleteCh() <- char
}

func (p *CharacterPipeline) Update(char Character) {
	p.charManagerStage.GetUpdateCh() <- char
}

func (p *CharacterPipeline) Read(id string) Character {
	intVal, err := strconv.Atoi(id)
	if err != nil {
		fmt.Printf("Occurred error when converting character id.")
	}
	return p.charManagerStage.Read(intVal)
}
