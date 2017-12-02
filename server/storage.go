package main

import (
	"container/heap"
	"sync"
	"time"
)

type entry struct {
	value  interface{}
	expiry *expiry
}

type storage struct {
	mutex    sync.RWMutex
	entries  map[string]entry
	expiries expiryPriorityQueue
}

func newStorage() *storage {
	e := make(map[string]entry)
	pq := make(expiryPriorityQueue, 0, 100)
	heap.Init(&pq)
	return &storage{sync.RWMutex{}, e, pq}
}

func (s *storage) set(k string, v interface{}) bool {
	if v == nil {
		s.del(k)
		return false
	}

	s.mutex.Lock()
	if e, ok := s.entries[k]; ok && e.expiry != nil {
		s.expiries.del(e.expiry)
	}
	s.entries[k] = entry{value: v}
	s.mutex.Unlock()

	return true
}

func (s *storage) get(k string) interface{} {
	if e := s.getEntry(k); e != nil {
		return e.value
	}
	return nil
}

func (s *storage) getEntry(k string) *entry {
	s.expireIfNeeded(k)

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if e, ok := s.entries[k]; ok {
		return &e
	}
	return nil
}

func (s *storage) del(k string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if e, ok := s.entries[k]; ok {
		if e.expiry != nil {
			s.expiries.del(e.expiry)
		}
		delete(s.entries, k)
	}
}

func (s *storage) exists(k string) bool {
	return s.get(k) != nil
}

func (s *storage) setExpire(k string, ttl int64) {
	e := s.getEntry(k)
	if e == nil {
		return
	}

	s.mutex.Lock()
	{
		if e.expiry != nil {
			s.expiries.update(e.expiry, ttl)
		} else {
			e.expiry = newExpiry(k, ttl)
			s.expiries.add(e.expiry)
			s.entries[k] = *e
		}
	}
	s.mutex.Unlock()

	s.expireIfNeeded(k)
}

func (s *storage) expireIfNeeded(k string) {
	s.mutex.Lock()
	if e, ok := s.entries[k]; ok && e.expiry != nil && e.expiry.ttl < time.Now().Unix() {
		s.expiries.del(e.expiry)
		delete(s.entries, k)
	}
	s.mutex.Unlock()
}

func (s *storage) len() (entries int, expires int) {
	return len(s.entries), len(s.expiries)
}

func (s *storage) removeExpired() {
	now := time.Now().Unix()

	s.mutex.Lock()
	defer s.mutex.Unlock()

	for s.expiries.Len() > 0 {
		if s.expiries[0].ttl > now {
			break
		}

		expiry := heap.Pop(&s.expiries).(*expiry)
		if _, ok := s.entries[expiry.key]; ok {
			delete(s.entries, expiry.key)
		}
	}
}
