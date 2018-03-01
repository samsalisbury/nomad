package nomad

import "sync"

type StartStopper struct {
	stopped   bool
	stoppedCh chan struct{}
	runningCh chan struct{}
	sync.RWMutex
}

func NewStartStopper() *StartStopper {
	s := &StartStopper{}
	go func() {
		for {
			s.Lock()
			if s.stopped {
				s.Unlock()
				select {
				default:
				case s.stoppedCh <- struct{}{}:
				}
			} else {
				s.Unlock()
				select {
				default:
					s.runningCh <- struct{}{}
				}
			}
		}
	}()
	return s
}

func (s *StartStopper) Stop() {
	s.Lock()
	defer s.Unlock()
	s.stopped = true
}

func (s *StartStopper) Start() {
	s.Lock()
	defer s.Unlock()
	s.stopped = false
}

func (s *StartStopper) Stopped() <-chan struct{} {
	return s.stoppedCh
}

func (s *StartStopper) IsStopped() bool {
	s.RLock()
	defer s.RUnlock()
	return s.stopped
}

func (s *StartStopper) Started() <-chan struct{} {
	return s.runningCh
}

func (s *StartStopper) WaitForStart() {
	<-s.runningCh
}
