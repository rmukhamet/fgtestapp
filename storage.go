package main

import (
	"fmt"
	"sync"
	"time"
)

type userLog struct {
	IPv4 map[string][]time.Time
}

type storage struct {
	mu    sync.RWMutex
	items map[int64]userLog
}

func newStorage() *storage {
	items := make(map[int64]userLog)
	return &storage{items: items}
}

func newUserLog(data *logData) (log userLog, err error) {
	log.IPv4 = make(map[string][]time.Time)
	ts, err := time.Parse("2006-01-02 15:04:05", data.TS)
	if err != nil {
		err = fmt.Errorf("error parsing time:  %s. error: %s", data.TS, err)
	}
	log.IPv4[data.IPv4] = []time.Time{ts}
	return
}
func (s *storage) addLog(data *logData) error {
	if !data.valid() {
		return fmt.Errorf("log data not valid")
	}
	s.mu.Lock()
	log, err := newUserLog(data)
	if err != nil {
		s.mu.Unlock()
		return fmt.Errorf("error in log format: %s", err)
	}
	if logE, ok := s.items[data.ID]; ok {
		if ipLog, ok := logE.IPv4[data.IPv4]; ok {
			ipLog = append(ipLog, log.IPv4[data.IPv4]...)
			logE.IPv4[data.IPv4] = ipLog
			s.items[data.ID] = logE
		} else {
			logE.IPv4[data.IPv4] = log.IPv4[data.IPv4]
			s.items[data.ID] = logE
		}
	} else {
		s.items[data.ID] = log
	}
	s.mu.Unlock()

	return nil
}
func (s *storage) getByUserID(userID int64) (log userLog, err error) {
	s.mu.RLock()
	log, found := s.items[userID]
	if !found {
		s.mu.RUnlock()
		err = fmt.Errorf("logs belong userID: %d not find", userID)
		return
	}
	s.mu.RUnlock()
	return
}
