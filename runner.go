package plumber

type Runner interface {
	System() System
	SetSystem(System)
}
type runner struct {
	sys System
}

func (r *runner) System() System {
	return r.sys
}
func (r *runner) SetSystem(s System) {
	r.sys = s
}

func NewRunner() Runner {
	return &runner{}
}
