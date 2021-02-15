package loadtest

type FinishTestStrategy interface {
	getFinishChan() chan bool
	finishTesting()
}

type OneChanFinishTestStrategy struct {
	doneChan   chan bool
	concurrent uint32
}

func GetOneChanStrategy(concurrent uint32) FinishTestStrategy {
	return &OneChanFinishTestStrategy{
		doneChan:   make(chan bool, concurrent),
		concurrent: concurrent,
	}
}

func (s *OneChanFinishTestStrategy) getFinishChan() chan bool {
	return s.doneChan
}

func (s *OneChanFinishTestStrategy) finishTesting() {
	for i := uint32(0); i < s.concurrent; i++ {
		s.doneChan <- true
	}
}

type MultipleChansFinishTestStrategy struct {
	doneChans []chan bool
}

func GetMultipleChansStrategy(concurrent uint32) FinishTestStrategy {
	return &MultipleChansFinishTestStrategy{
		doneChans: []chan bool{},
	}
}

func (s *MultipleChansFinishTestStrategy) getFinishChan() chan bool {
	doneChan := make(chan bool, 1)
	s.doneChans = append(s.doneChans, doneChan)

	return doneChan
}

func (s *MultipleChansFinishTestStrategy) finishTesting() {
	for _, doneChan := range s.doneChans {
		doneChan <- true
	}
}
