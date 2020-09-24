package plumber

// System
type System interface {
	Name() string
	State() State
	SetState(State) System
	Map(Lambda) System
	From(Stream) System
	To(Stream) System
	Initiate() chan error
}
type SystemConfigurer func(s System) System
type system struct {
	name  string
	state State
	nodes []Lambda
	in    Stream
	out   Stream
}

func (s *system) Name() string {
	return s.name
}

func (s *system) State() State {
	return s.state
}

func (s *system) SetState(st State) System {
	s.state = st
	return s
}

func (s *system) Map(l Lambda) System {
	s.nodes = append(s.nodes, l)
	return s
}

func (s *system) From(st Stream) System {
	s.in = st
	return s
}

func (s *system) To(st Stream) System {
	s.out = st
	return s
}

func (s *system) Initiate() chan error {
	var lcs []*lamdaContainer
	for idx, n := range s.nodes {
		lc := &lamdaContainer{}
		lc.l = n
		if idx == 0 {
			lc.In = s.in
		} else {
			lc.In = lcs[idx-1].Out
		}
		if idx == len(s.nodes)-1 {
			lc.Out = s.out
		} else {
			lc.Out = make(chan interface{})
		}
		lcs = append(lcs, lc)
	}
	errs := make(chan error, 1024) //TODO: configure error chan cap
	for _, lc := range lcs {
		go func(container *lamdaContainer) {
			for v := range container.In {
				v, err := container.l(s.State(), v)
				if err != nil {
					errs <- err
					continue
				}
				container.Out <- v
			}
		}(lc)
	}
	return errs
}

func setDefaultSystemConfigs(s System) System {

	return s
}
func NewDefaultSystem(confs ...SystemConfigurer) System {
	confs = append(confs, setDefaultSystemConfigs)
	var s System
	s = &system{}
	for _, c := range confs {
		s = c(s)
	}
	return s
}
