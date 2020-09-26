package system

import (
	"github.com/amirrezaask/plumber"
	"github.com/amirrezaask/plumber/stream"
)

type lamdaContainer struct {
	l   plumber.Lambda
	In  plumber.Stream
	Out plumber.Stream
}
type SystemConfigurer func(s plumber.System) plumber.System

type defaultSystem struct {
	name  string
	state plumber.State
	nodes []plumber.Lambda
	in    plumber.Stream
	out   plumber.Stream
}

func (s *defaultSystem) Name() string {
	return s.name
}

func (s *defaultSystem) State() plumber.State {
	return s.state
}

func (s *defaultSystem) SetState(st plumber.State) plumber.System {
	s.state = st
	return s
}

func (s *defaultSystem) Then(l plumber.Lambda) plumber.System {
	s.nodes = append(s.nodes, l)
	return s
}

func (s *defaultSystem) From(st plumber.Stream) plumber.System {
	s.in = st
	return s
}

func (s *defaultSystem) To(st plumber.Stream) plumber.System {
	s.out = st
	return s
}

func (s *defaultSystem) Initiate() (chan error, error) {
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
			lc.Out = stream.NewDumbStream()
		}
		lcs = append(lcs, lc)
	}
	errs := make(chan error, 1024) //TODO: configure error chan cap
	for _, lc := range lcs {
		err := lc.In.StartReading()
		if err != nil {
			return nil, err
		}
		go func(container *lamdaContainer) {
			for v := range container.In.ReadChan() {
				f := func(i interface{}) {
					i = v
				}
				v, err := container.l(s.State(), f)
				if err != nil {
					errs <- err
					continue
				}
				container.Out.Write(v)
			}
		}(lc)
	}
	return errs, nil
}

func setDefaultSystemConfigs(s plumber.System) plumber.System {
	return s
}
func NewDefaultSystem(confs ...SystemConfigurer) plumber.System {
	confs = append(confs, setDefaultSystemConfigs)
	var s plumber.System
	s = &defaultSystem{}
	for _, c := range confs {
		s = c(s)
	}
	return s
}
