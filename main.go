package main

import (
	"time"

	character "github.com/maxabcr2000/ChannelPractice/character"
)

func main() {
	outBoundCh := make(chan character.Character)

	sinkStage := character.NewSinkStage(outBoundCh)
	charManagerStage := character.NewCharManagementStage(outBoundCh)
	sourceStage := character.NewSourceStage(charManagerStage.GetRegisterCh())

	go sourceStage.Start()
	go charManagerStage.Start()
	go sinkStage.Start()

	time.Sleep(50 * time.Second)
}
