package psetter

import (
	"time"

	"github.com/nickwells/check.mod/check"
)

// Time Formats given here can be used to set the Format member of the Time
// struct
const (
	TimeFmtDefault   = "2006/Jan/02T15:04:05"
	TimeFmtHMS       = "15:04:05"
	TimeFmtHoursMins = "15:04"
	TimeFmtDateOnly  = "2006/Jan/02"
	TimeFmtTimestamp = "20060102.150405"
	TimeFmtISO8601   = "2006-01-02T15:04:05"
)

// Time allows you to give a parameter that can be used to set a time.Time
// value.
type Time struct {
	ValueReqMandatory

	// You must set a Value, the program will panic if not. This is the
	// Time that the setter is setting.
	Value *time.Time
	// The Format is used to convert the string parameter value into a Time
	// suitable for setting the Value. If no Format is given the default
	// value will be used, see TimeFmtDefault.
	Format string
	// The Checks, if any, are applied to the time.Time value and the new
	// parameter will be applied only if they all return a nil error.
	Checks []check.Time
}

// CountChecks returns the number of check functions this setter has
func (s Time) CountChecks() int {
	return len(s.Checks)
}

// SetWithVal (called when a value follows the parameter) checks that the
// value can be parsed to a time, if it cannot be parsed successfully it
// returns an error. If there is a check and the check is violated it returns
// an error. Only if the value is parsed successfully and the check is not
// violated is the Value set.
func (s Time) SetWithVal(_ string, paramVal string) error {
	v, err := time.Parse(s.format(), paramVal)
	if err != nil {
		return err
	}

	for _, check := range s.Checks {
		if check == nil {
			continue
		}

		err := check(v)
		if err != nil {
			return err
		}
	}

	*s.Value = v
	return nil
}

// format returns the format string if set or else the default value
func (s Time) format() string {
	if s.Format != "" {
		return s.Format
	}
	return TimeFmtDefault
}

// AllowedValues returns a string describing the allowed values
func (s Time) AllowedValues() string {
	return "any value that represents a time that can be parsed using the" +
		" time format string: " + s.format() +
		HasChecks(s)
}

// CurrentValue returns the current setting of the parameter value
func (s Time) CurrentValue() string {
	return s.Value.String()
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil.
func (s Time) CheckSetter(name string) {
	if s.Value == nil {
		panic(NilValueMessage(name, "psetter.Time"))
	}
}
