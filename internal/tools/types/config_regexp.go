package types

import (
	"regexp"
)

type Regexp struct {
	*regexp.Regexp
}

func (r *Regexp) FromString(s string) error {
	re, err := regexp.Compile(s)
	if err != nil {
		return err
	}
	r.Regexp = re
	return nil
}
