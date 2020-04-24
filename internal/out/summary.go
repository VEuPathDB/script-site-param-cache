package out

import (
	"encoding/json"
	_ "gopkg.in/yaml.v3"
	"sync"
	"time"
)

type Details struct {
	Succeeded uint `json:"pass" yaml:"pass"`
	Failed    uint `json:"fail" yaml:"fail"`
	Skipped   uint `json:"skip" yaml:"skip"`
}

type Summary struct {
	Url           string                   `json:"url" yaml:"URL"`
	RecordTypes   *Details                 `json:"recordTypes,omitempty" yaml:"RecordTypes,omitempty"`
	SearchDetails *Details                 `json:"searchDetails,omitempty" yaml:"SearchDetails,omitempty"`
	Timings       map[string]time.Duration `json:"-" yaml:"-"`
	lock          sync.Mutex
}

func (s *Summary) ensureRt() {
	if s.RecordTypes == nil {
		s.RecordTypes = &Details{}
	}
}

func (s *Summary) RecordTypeSuccess() {
	s.lock.Lock()
	s.ensureRt()
	s.RecordTypes.Succeeded++
	s.lock.Unlock()
}

func (s *Summary) RecordTypeFailed() {
	s.lock.Lock()
	s.ensureRt()
	s.RecordTypes.Failed++
	s.lock.Unlock()
}

func (s *Summary) RecordTypeSkipped() {
	s.lock.Lock()
	s.ensureRt()
	s.RecordTypes.Skipped++
	s.lock.Unlock()
}

func (s *Summary) ensureSd() {
	if s.SearchDetails == nil {
		s.SearchDetails = &Details{}
	}
}

func (s *Summary) SearchDetailSuccess() {
	s.lock.Lock()
	s.ensureSd()
	s.SearchDetails.Succeeded++
	s.lock.Unlock()
}

func (s *Summary) SearchDetailFailed() {
	s.lock.Lock()
	s.ensureSd()
	s.SearchDetails.Failed++
	s.lock.Unlock()
}

func (s *Summary) SearchDetailSkipped() {
	s.lock.Lock()
	s.ensureSd()
	s.SearchDetails.Skipped++
	s.lock.Unlock()
}

func (s *Summary) RecordTiming(url string, dur time.Duration) {
	s.lock.Lock()
	if s.Timings == nil {
		s.Timings = make(map[string]time.Duration)
	}
	s.Timings[url] = dur
	s.lock.Unlock()
}

func (s Summary) MarshalJSON() ([]byte, error) {
	type alias Summary
	return json.Marshal(struct {
		Summary alias `json:"summary"`
	}{alias(s)})
}

func (s Summary) MarshalYAML() (interface{}, error) {
	type alias Summary
	return struct {
		Summary alias `yaml:"Summary"`
	}{alias(s)}, nil
}
