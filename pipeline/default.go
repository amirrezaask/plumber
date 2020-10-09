package pipeline

import (
	"github.com/amirrezaask/plumber"
)

type container struct {
	pipeCtx *plumber.PipeCtx
	pipe    plumber.Pipe
}
type SystemConfigurer func(s plumber.Pipeline) plumber.Pipeline

type defaultSystem struct {
	name       string
	errs       chan error
	checkpoint plumber.Checkpoint
	state      plumber.State
	nodes      []plumber.Pipe
	in         plumber.Stream
	out        plumber.Stream
}

func (s *defaultSystem) Errors() chan error {
	return s.errs
}
func (s *defaultSystem) Checkpoint() {
	s.checkpoint(s)
}

func (s *defaultSystem) SetCheckpoint(c plumber.Checkpoint) plumber.Pipeline {
	s.checkpoint = c
	return s
}
func (s *defaultSystem) InputStream() plumber.Stream {
	return s.in
}
func (s *defaultSystem) OutputStream() plumber.Stream {
	return s.out
}
func (s *defaultSystem) Name() string {
	return s.name
}

func (s *defaultSystem) State() plumber.State {
	return s.state
}

func (s *defaultSystem) SetState(st plumber.State) plumber.Pipeline {
	s.state = st
	return s
}

func (s *defaultSystem) Then(l plumber.Pipe) plumber.Pipeline {
	s.nodes = append(s.nodes, l)
	return s
}
func (s *defaultSystem) Thens(ls ...plumber.Pipe) plumber.Pipeline {
	s.nodes = append(s.nodes, ls...)
	return s
}

func (s *defaultSystem) From(st plumber.Stream) plumber.Pipeline {
	s.in = st
	return s
}

func (s *defaultSystem) UpdateState() error {
	err := s.State().Set(s.in.Name(), s.in.State())
	if err != nil {
		return err
	}
	err = s.State().Set(s.out.Name(), s.out.State())
	if err != nil {
		return err
	}
	return nil
}

func (s *defaultSystem) To(st plumber.Stream) plumber.Pipeline {
	s.out = st
	return s
}

func (s *defaultSystem) Initiate() (chan error, error) {
	// update streams state if there is any checkpoints
	inState, err := s.State().Get(s.in.Name())
	if err != nil {
		return nil, err
	}

	outState, err := s.State().Get(s.out.Name())
	if err != nil {
		return nil, err
	}
	if inState != nil {
		s.in.LoadState(inState.(map[string]interface{}))
	}
	if outState != nil {
		s.out.LoadState(outState.(map[string]interface{}))
	}
	// connect lambdas through channels
	var lcs []*container
	if err != nil {
		return nil, err
	}
	for idx, n := range s.nodes {
		lc := &container{
			pipeCtx: &plumber.PipeCtx{},
		}
		lc.pipe = n
		if idx == 0 {
			lc.pipeCtx.In = s.in.Input()
		} else {
			lc.pipeCtx.In = lcs[idx-1].pipeCtx.Out
		}
		if idx == len(s.nodes)-1 {
			lc.pipeCtx.Out = s.OutputStream().Output()
		} else {
			lc.pipeCtx.Out = make(chan interface{})
		}
		lcs = append(lcs, lc)
	}
	errs := make(chan error, 1024) //TODO: configure error chan cap

	//start checkpoint process
	go s.Checkpoint()
	for _, lc := range lcs {
		go func(container *container) {
			container.pipeCtx.Err = errs
			container.pipeCtx.State = s.State()
			container.pipe(container.pipeCtx)
		}(lc)
	}
	return errs, nil
}

func setDefaultSystemConfigs(s plumber.Pipeline) plumber.Pipeline {
	return s
}
func NewDefaultSystem(confs ...SystemConfigurer) plumber.Pipeline {
	confs = append(confs, setDefaultSystemConfigs)
	var s plumber.Pipeline
	s = &defaultSystem{}
	for _, c := range confs {
		s = c(s)
	}
	return s
}
