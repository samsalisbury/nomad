package nomad

import "sync"

type StartStopper struct {
	stoppedCh chan struct{}
	sync.RWMutex
}

func NewStartStopper() *StartStopper {
	return &StartStopper{stoppedCh: make(chan struct{})}
}

func (s *StartStopper) Stop() {
	s.Lock()
	defer s.Unlock()
	close(s.stoppedCh)
}

func (s *StartStopper) Start() {
	s.Lock()
	defer s.Unlock()
	select {
	default:
		// Already started.
	case <-s.stoppedCh:
		s.stoppedCh = make(chan struct{})
	}
}

func (s *StartStopper) Stopped() <-chan struct{} {
	s.RLock()
	defer s.RUnlock()
	return s.stoppedCh
}

func (s *StartStopper) IsStopped() bool {
	s.RLock()
	defer s.RUnlock()
	select {
	default:
		return false
	case <-s.stoppedCh:
		return true
	}
}
