package gowhere

import (
	"fmt"
)

// InvalidCond represents the error when invalid condition is given
type InvalidCond struct {
	// The given condition(s)
	cond interface{}
	// The vars if given
	vars interface{}
}

func (e *InvalidCond) Error() string {
	return fmt.Sprintf("Invalid Conditions: %+v %+v", e.cond, e.vars)
}
