// model.go
package enum

import (
	"fmt"
)

type NumberRange struct {
	Upper   uint64   `json:"upper"`
	Lower   uint64   `json:"lower"`
	Records []Record `json:"records"`
}

type Record struct {
	Order       uint16 `json:"order"`
	Preference  uint16 `json:"preference"`
	Flags       string `json:"flags"`
	Service     string `json:"service"`
	Regexp      string `json:"regexp"`
	Replacement string `json:"replacement"`
}

func (r *NumberRange) OverlapWith(o NumberRange) bool {
	if o.Lower <= r.Lower && r.Upper <= o.Upper {
		return true
	} else {
		return false
	}
}

// RangeOverlapError is returned when an operation fails because
// a range overlaps with on or more other ranges.
type RangeOverlapError struct {
	Range    NumberRange
	Overlaps []NumberRange
}

func (e *RangeOverlapError) Error() string {
	if len(e.Overlaps) == 1 {
		return fmt.Sprintf("[%15.d:%15.d] orverlaps with [%15.d:%15.d]",
			e.Range.Lower, e.Range.Upper,
			e.Overlaps[0].Upper, e.Overlaps[0].Upper)
	} else {
		return fmt.Sprintf("[%15.d:%15.d] orverlaps with %d other ranges", e.Range.Lower, e.Range.Upper, len(e.Overlaps))
	}
}

type Backend interface {
	// RangesBetween returns a list of ranges that enclose the given range l(ower) to u(pper) or
	// nil if no range matches.
	// The c parameter is the maximum count of values to return. If a negative c value is used
	// it will return the ranges in reverse order.
	RangesBetween(l, u uint64, c int) ([]NumberRange, error)

	// Add a range to the backend. Any range overlapping with the one added will be deleted or
	// adjusted to make room for the new one and returned.
	PushRange(r NumberRange) ([]NumberRange, error)

	// Close the backend.
	Close() error
}
