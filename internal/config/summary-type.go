package config

import "fmt"

const (
	errBadSumType = `unrecognized summary type "%s".  Valid types are %s, %s.`
)

type SummaryType string

const (
	SummaryTypeJson SummaryType = "json"
	SummaryTypeYaml SummaryType = "yaml"
)

func (s *SummaryType) Unmarshal(value string) (err error) {
	tmp := SummaryType(value)
	switch tmp {
	case SummaryTypeJson:
		fallthrough
	case SummaryTypeYaml:
		*s = tmp
	default:
		return fmt.Errorf(errBadSumType, value, SummaryTypeJson, SummaryTypeYaml)
	}
	return nil
}
