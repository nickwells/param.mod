package psetter

import (
	"fmt"
	"strings"
	"time"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/english.mod/english"
	"github.com/nickwells/strdist.mod/v2/strdist"
)

// TimeLocation allows you to give a parameter that can be used to set a
// time.Location pointer. You can also supply check functions that will
// validate the Value.
type TimeLocation struct {
	ValueReqMandatory

	// You must set a Value, the program will panic if not. Note that this is
	// a pointer to the pointer to the Location, you should initialise it
	// with the address of the Location pointer.
	Value **time.Location

	// Locations (which can be empty) is used to provide an improved error
	// message when the location cannot be loaded.
	Locations []string

	// The Checks, if any, are applied to the supplied parameter value and
	// the new parameter will be applied only if they all return a nil error.
	Checks []check.TimeLocation
}

// CountChecks returns the number of check functions this setter has
func (s TimeLocation) CountChecks() int {
	return len(s.Checks)
}

// makeSuggestionStr formats the alternate locations into a message suitable
// for presentation to the user.
func (s TimeLocation) makeSuggestionStr(altLocs []string) string {
	if len(altLocs) == 0 {
		return ""
	}

	preamble := ", did you mean "
	if len(altLocs) > 1 {
		preamble += "one of "
	}

	return preamble + english.JoinQuoted(altLocs, ", ", " or ", `"`, `"`)
}

// suggestAltTimeLocation tries to find values in the list of available
// locations which are similar to the badLoc value
func (s TimeLocation) suggestAltTimeLocation(badLoc string) string {
	const alternativeCount = 3

	if len(s.Locations) == 0 {
		return ""
	}

	finder := strdist.DefaultFinders[strdist.CaseBlindAlgoNameCosine]

	var altLocs []string
	for _, f := range []func(string, []string) []string{
		// This finds matches against the locations
		func(s string, locs []string) []string {
			return finder.FindNStrLike(alternativeCount, s, locs...)
		},
		// This finds those entries in Locations which have a name with parts
		// separated by a '/'. This is how geographical timezone locations
		// are represented, for instance Europe/London, Asia/Jerusalem or
		// America/Indiana/Indianapolis. The match is just against the last
		// part of the name, the city name (London, Jerusalem or
		// Indianapolis). Any discovered matches are then mapped back to the
		// original full name.
		func(s string, locs []string) []string {
			justCities := []string{} // a misnomer as they aren't all city names
			backMap := map[string][]string{}

			for _, l := range locs {
				parts := strings.Split(l, "/")
				if len(parts) > 1 {
					city := parts[len(parts)-1]
					justCities = append(justCities, city)
					backMap[city] = append(backMap[city], l)
				}
			}

			matches := finder.FindNStrLike(alternativeCount, s, justCities...)
			if len(matches) == 0 {
				return matches
			}

			rval := []string{}
			for _, m := range matches {
				rval = append(rval, backMap[m]...)
			}

			return rval
		},
	} {
		altLocs = f(badLoc, s.Locations)
		if len(altLocs) > 0 {
			break
		}

		altLocs = f(strings.ReplaceAll(badLoc, " ", "_"), s.Locations)
		if len(altLocs) > 0 {
			break
		}
	}

	return s.makeSuggestionStr(altLocs)
}

// SetWithVal (called when a value follows the parameter) checks that the
// value can be parsed to a Location (note that it will try to replace any
// spaces with underscores if the first attempt fails). If it cannot be
// parsed successfully it returns an error. The Checks, if any, will be
// applied and if any of them return a non-nil error the Value will not be
// updated and the error will be returned.
func (s TimeLocation) SetWithVal(_ string, paramVal string) error {
	v, err := time.LoadLocation(paramVal)
	if err != nil {
		convertedVal := strings.ReplaceAll(paramVal, " ", "_")

		var e2 error

		v, e2 = time.LoadLocation(convertedVal)
		if e2 != nil {
			return fmt.Errorf("bad timezone %q%s: %w",
				paramVal, s.suggestAltTimeLocation(paramVal), err)
		}
	}

	for _, check := range s.Checks {
		err := check(*v)
		if err != nil {
			return err
		}
	}

	*s.Value = v

	return nil
}

// AllowedValues returns a string describing the allowed values
func (s TimeLocation) AllowedValues() string {
	return "any value that represents a location" +
		HasChecks(s) +
		". Typically this will be a string of the form" +
		" Continent/City_Name, for instance, Europe/London" +
		" or America/New_York." +
		" Additionally some of the three-letter timezone" +
		" names are also allowed such as UTC or CET."
}

// CurrentValue returns the current setting of the parameter value
func (s TimeLocation) CurrentValue() string {
	return (*s.Value).String()
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil or if it has nil Checks.
func (s TimeLocation) CheckSetter(name string) {
	// Check the value is not nil
	if s.Value == nil {
		panic(NilValueMessage(name, fmt.Sprintf("%T", s)))
	}

	// Check there are no nil Check funcs
	for i, check := range s.Checks {
		if check == nil {
			panic(NilCheckMessage(name, fmt.Sprintf("%T", s), i))
		}
	}
}
