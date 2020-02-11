package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/bxcodec/faker/v3"
)

func TestCheckDupes(t *testing.T) {
	var (
		err        error
		pathGolden = filepath.Join("testdata", "testlog.golden")
		logs       []logData
	)
	testData, err := ioutil.ReadFile(pathGolden)
	if err != nil {
		t.Fatalf("failed reading .golden: %s", err)
	}
	err = json.Unmarshal(testData, &logs)
	if err != nil {
		t.Error(err)
	}

	type test struct {
		userID1 int64
		userID2 int64
		want    bool
	}

	tests := []test{
		{userID1: 1, userID2: 2, want: true},
		{userID1: 1, userID2: 3, want: false},
		{userID1: 2, userID2: 1, want: true},
		{userID1: 2, userID2: 3, want: true},
		{userID1: 3, userID2: 2, want: true},
		{userID1: 1, userID2: 4, want: false},
		{userID1: 3, userID2: 1, want: false},
		{userID1: 1, userID2: 1, want: true},
	}
	s := newStorage()
	for _, log := range logs {
		err = s.addLog(&log)
		if err != nil {
			t.Error(err)
		}
	}
	ls := newLogService(s)
	for _, test := range tests {
		r, err := ls.checkDupes(test.userID1, test.userID2)
		if err != nil {
			t.Error(err)
		}
		if r != test.want {
			t.Errorf("userID1: %d userID2: %d. Func result: %v expected: %v", test.userID1, test.userID2, r, test.want)
		}

	}

}

func BenchmarkCheckDupes(b *testing.B) {

	var (
		err        error
		pathGolden = filepath.Join("testdata", "testlogbig.golden")
	)
	if _, err := os.Stat(pathGolden); err != nil {
		err = generateFakeLog(pathGolden)
		if err != nil {
			b.Error(err)
			return
		}
	}

	s := newStorage()
	ls := newLogService(s)
	err = ls.getData(pathGolden)
	if err != nil {
		b.Error(err)
		return
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ls.checkDupes(5, 60)
	}
}
func generateFakeLog(path string) error {
	type logTest struct {
		ID   int64  `faker:"boundary_start=1, boundary_end=100000"`
		IPv4 string `faker:"ipv4"`
		TS   string `faker:"timestamp"`
	}
	const limit = 15000000
	logs := make([]logTest, limit)
	for i := 0; i < limit; i++ {
		l := logTest{}
		err := faker.FakeData(&l)
		if err != nil {
			err = fmt.Errorf("failed generating fake data: %s", err)
			return err
		}
		logs[i] = l
	}
	logsB, err := json.Marshal(logs)
	err = ioutil.WriteFile(path, logsB, 0644)
	if err != nil {
		err = fmt.Errorf("failed writing .golden: %s", err)
		return err
	}
	return nil
}
