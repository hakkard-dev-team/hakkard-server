package commands



func NewNameMatcher(r *Route) func(string) bool {
	return func(cmd string) bool {
		return cmd == r.Name
	}
}
