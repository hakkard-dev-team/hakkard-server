package game

import (
	"fmt"

	"github.com/Matt-Gleich/logoru"
)



func NewNameMatcher(r *Route) func(string) bool {
	return func(cmd string) bool {
		logoru.Debug(fmt.Sprintf("Attempting to match %s to %s", r.Name, cmd))
		return cmd == r.Name
	}
}
