package util

import "time"

type RequestTimeout time.Duration

func (r *RequestTimeout) UnmarshalFlag(value string) error {
	tmp, err := time.ParseDuration(value)
	*r = RequestTimeout(tmp)
	return err
}

func (r RequestTimeout) ToDuration() time.Duration {
	return time.Duration(r)
}
