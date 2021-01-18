package models

import "sync/atomic"

type ＷorkerPool struct {
	workers []*Worker
	current uint64
}

func (s *ＷorkerPool) NextIndex() int {
	return int(atomic.AddUint64(&s.current, uint64(1)) % uint64(len(s.workers)))
}

// GetNextPeer returns next active peer to take a connection
func (s *ＷorkerPool) GetNextPeer() *Ｗorker {
	// loop entire backends to find out an Alive backend
	next := s.NextIndex()
	l := len(s.workers) + next // start from next and move a full cycle
	for i := next; i < l; i++ {
		idx := i % len(s.workers) // take an index by modding with length
		// if we have an alive backend, use it and store if its not the original one
		if s.workers[idx].IsAlive() {
			if i != next {
				atomic.StoreUint64(&s.current, uint64(idx)) // mark the current one
			}
			return s.workers[idx]
		}
	}
	return nil
}
