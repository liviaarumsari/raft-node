package raft

import (
	"if3230-tubes-wreckitraft/shared"
	"sync"
	"sync/atomic"
)

type raftState struct {
	currentTerm    uint64 // cache of stable store
	commitIndex    uint64
	lastAppliedLog uint64
	lastLogIndex   uint64 // cache of log store
	lastLogTerm    uint64
	nextIndex      map[shared.Address]uint64
	matchIndex     map[shared.Address]uint64
	routinesGroup  sync.WaitGroup
	state          NodeType
	lock           sync.Mutex
}

func (r *raftState) getState() NodeType {
	stateAddr := (*uint32)(&r.state)
	return NodeType(atomic.LoadUint32(stateAddr))
}

func (r *raftState) setState(s NodeType) {
	stateAddr := (*uint32)(&r.state)
	atomic.StoreUint32(stateAddr, uint32(s))
}

func (r *raftState) getCurrentTerm() uint64 {
	return atomic.LoadUint64(&r.currentTerm)
}

func (r *raftState) setCurrentTerm(term uint64) {
	atomic.StoreUint64(&r.currentTerm, term)
}

func (r *raftState) getCommitIndex() uint64 {
	return atomic.LoadUint64(&r.commitIndex)
}

func (r *raftState) setCommitIndex(commitIndex uint64) {
	atomic.StoreUint64(&r.commitIndex, commitIndex)
}

func (r *raftState) getNextIndex(key shared.Address) (index uint64, bool2 bool) {
	r.lock.Lock()
	index, ok := r.nextIndex[key]
	r.lock.Unlock()
	return index, ok
}

func (r *raftState) setNextIndex(key shared.Address, index uint64) {
	r.lock.Lock()
	r.nextIndex[key] = index
	r.lock.Unlock()
}

func (r *raftState) getMatchIndex(key shared.Address) (index uint64, bool2 bool) {
	r.lock.Lock()
	index, ok := r.matchIndex[key]
	r.lock.Unlock()
	return index, ok
}

func (r *raftState) setMatchIndex(key shared.Address, index uint64) {
	r.lock.Lock()
	r.matchIndex[key] = index
	r.lock.Unlock()
}

func (r *raftState) getLastLog() (index, term uint64) {
	r.lock.Lock()
	index = r.lastLogIndex
	term = r.lastLogTerm
	r.lock.Unlock()
	return index, term
}

func (r *raftState) setLastLog(index, term uint64) {
	r.lock.Lock()
	r.lastLogIndex = index
	r.lastLogTerm = term
	r.lock.Unlock()
}

func (r *raftState) getLastIndex() uint64 {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.lastLogIndex
}

func (r *raftState) goFunc(f func()) {
	r.routinesGroup.Add(1)
	go func() {
		defer r.routinesGroup.Done()
		f()
	}()
}
