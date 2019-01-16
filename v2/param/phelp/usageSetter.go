package phelp

import (
	"github.com/nickwells/param.mod/v2/param"
)

// helpType records the helper to be updated and the style to apply if this
// parameter is set
type helpType struct {
	style helpStyle
	h     *StdHelp
}

func (helpType) ValueReq() param.ValueReq { return param.None }

func (s helpType) Set(_ string) error {
	s.h.showHelp = true
	s.h.style = s.style
	return nil
}

func (s helpType) SetWithVal(_ string, _ string) error {
	return s.Set("")
}

func (helpType) AllowedValues() string { return "none" }

func (helpType) CurrentValue() string { return "" }

func (helpType) CheckSetter(_ string) {}

// helpFull records the helper to be updated and will set the showHelp and
// showAllParams flags
type helpFull struct {
	h *StdHelp
}

func (helpFull) ValueReq() param.ValueReq { return param.None }

func (s helpFull) Set(_ string) error {
	s.h.showHelp = true
	s.h.showAllParams = true
	return nil
}

func (s helpFull) SetWithVal(_ string, _ string) error {
	return s.Set("")
}

func (helpFull) AllowedValues() string { return "none" }

func (helpFull) CurrentValue() string { return "" }

func (helpFull) CheckSetter(_ string) {}
