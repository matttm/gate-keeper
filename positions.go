package main

type RelativePosition int
type RelativePositionStr string

const (
	BEFORE RelativePosition = iota - 1
	INSIDE
	AFTER
)
const (
	BEFORE_ST = "Outside before"
	INSIDE_ST = "Inside"
	AFTER_ST  = "Outside after"
)

func getPositionOptions() []string {
	return []string{BEFORE_ST, INSIDE_ST, AFTER_ST}
}
func (p RelativePosition) String() RelativePositionStr {
	return RelativePositionStr(getPositionOptions()[p+1])
}
func (p RelativePositionStr) Value() RelativePosition {
	switch p {
	case BEFORE_ST:
		return BEFORE
	case INSIDE_ST:
		return INSIDE
	case AFTER_ST:
		return AFTER
	}
	return -100
}
