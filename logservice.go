package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type logData struct {
	ID   int64
	IPv4 string
	TS   string
}

type logService struct {
	storage *storage
}

func newLogService(s *storage) *logService {
	return &logService{storage: s}
}
func (ld *logData) valid() bool {
	return ld.TS != "" || ld.IPv4 != ""
}
func (ls *logService) checkDupes(userID1 int64, userID2 int64) (bool, error) {
	var matches int
	logsU1, err := ls.storage.getByUserID(userID1)
	if err != nil {
		return false, err
	}

	logsU2, err := ls.storage.getByUserID(userID2)
	if err != nil {
		return false, err
	}

	if len(logsU1.IPv4) > len(logsU2.IPv4) {
		for ip := range logsU2.IPv4 {
			if _, ok := logsU1.IPv4[ip]; ok {
				matches++
			}
		}
	} else {
		for ip := range logsU1.IPv4 {
			if _, ok := logsU2.IPv4[ip]; ok {
				matches++
			}
		}
	}
	if matches > 1 {
		return true, nil
	}

	return false, nil
}

func (ls *logService) getData(path string) (err error) {
	logsFile, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	logData := make([]logData, 100)
	err = json.Unmarshal(logsFile, &logData)
	if err != nil {
		return
	}
	if len(logData) == 0 {
		err = fmt.Errorf("no log data in file")
		return
	}
	for _, log := range logData {
		err = ls.storage.addLog(&log)
		if err != nil {
			return
		}
	}
	return
}
