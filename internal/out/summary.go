package out

import (
	"encoding/json"
	_ "gopkg.in/yaml.v3"
)

type Details struct {
	Succeeded uint `json:"pass" yaml:"pass"`
	Failed    uint `json:"fail" yaml:"fail"`
	Skipped   uint `json:"skip" yaml:"skip"`
}

type Summary struct {
	Url           string   `json:"url" yaml:"URL"`
	RecordTypes   *Details `json:"recordTypes,omitempty" yaml:"RecordTypes,omitempty"`
	SearchDetails *Details `json:"searchDetails,omitempty" yaml:"SearchDetails,omitempty"`
}

func (s *Summary) ensureRt() {
	if s.RecordTypes == nil {
		s.RecordTypes = &Details{}
	}
}

func (s *Summary) RecordTypeSuccess() {
	s.ensureRt()
	s.RecordTypes.Succeeded++
}

func (s *Summary) RecordTypeFailed() {
	s.ensureRt()
	s.RecordTypes.Failed++
}

func (s *Summary) RecordTypeSkipped() {
	s.ensureRt()
	s.RecordTypes.Skipped++
}

func (s *Summary) ensureSd() {
	if s.SearchDetails == nil {
		s.SearchDetails = &Details{}
	}
}

func (s *Summary) SearchDetailSuccess() {
	s.ensureSd()
	s.SearchDetails.Succeeded++
}

func (s *Summary) SearchDetailFailed() {
	s.ensureSd()
	s.SearchDetails.Failed++
}

func (s *Summary) SearchDetailSkipped() {
	s.ensureSd()
	s.SearchDetails.Skipped++
}

func (s Summary) MarshalJSON() ([]byte, error) {
	type alias Summary
	return json.Marshal(struct { Summary alias `json:"summary"` }{ alias(s) })
}

func (s Summary) MarshalYAML() (interface{}, error) {
	type alias Summary
	return struct {Summary alias `yaml:"Summary"`}{alias(s)}, nil
}

