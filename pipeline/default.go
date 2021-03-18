package pipeline

import (
	"bytes"
	"errors"
	"github.com/amirrezaask/plumber"
	"github.com/amirrezaask/plumber/pipe"
)

type SystemConfigurer func(s plumber.Pipeline) plumber.Pipeline

func pass(_ plumber.State, i interface{}) (interface{}, error) {
	return i, nil
}

type defaultPipeline struct {
	name       string
	errs       chan error
	checkpoint plumber.Checkpoint
	state      plumber.State
	nodes      []plumber.Pipe
	in         plumber.Input
	out        plumber.Output
	logger     plumber.Logger
}

func (s *defaultPipeline) Errors() chan error {
	return s.errs
}
func (s *defaultPipeline) Checkpoint() {
	s.checkpoint(s)
}

func (s *defaultPipeline) SetCheckpoint(c plumber.Checkpoint) plumber.Pipeline {
	s.checkpoint = c
	return s
}
func (s *defaultPipeline) GetInputStream() plumber.Input {
	return s.in
}
func (s *defaultPipeline) GetOutputStream() plumber.Output {
	return s.out
}
func (s *defaultPipeline) Name() string {
	return s.name
}

func (s *defaultPipeline) State() plumber.State {
	return s.state
}

func (s *defaultPipeline) SetState(st plumber.State) plumber.Pipeline {
	s.state = st
	return s
}

func (s *defaultPipeline) Then(l plumber.Pipe) plumber.Pipeline {
	s.nodes = append(s.nodes, l)
	return s
}
func (s *defaultPipeline) Thens(ls ...plumber.Pipe) plumber.Pipeline {
	s.nodes = append(s.nodes, ls...)
	return s
}

func (s *defaultPipeline) From(st plumber.Input) plumber.Pipeline {
	s.in = st
	return s
}

func (s *defaultPipeline) UpdateState() error {
	inState, err := s.in.State()
	if err != nil {
		return err
	}
	err = s.State().Set(s.in.Name(), inState)
	if err != nil {
		return err
	}
	outState, err := s.out.State()
	if err != nil {
		return err
	}
	err = s.State().Set(s.out.Name(), outState)
	if err != nil {
		return err
	}
	return nil
}

func (s *defaultPipeline) To(st plumber.Output) plumber.Pipeline {
	s.out = st
	return s
}

func (s *defaultPipeline) WithLogger(l plumber.Logger) plumber.Pipeline {
	s.logger = l
	return s
}

func (s *defaultPipeline) Logger() plumber.Logger {
	return s.logger
}

func (s *defaultPipeline) Initiate() (chan error, error) {
	// update streams state if there is any checkpoints
	inState, err := s.State().GetBytes(s.in.Name())
	if err != nil {
		return nil, err
	}

	outState, err := s.State().GetBytes(s.out.Name())
	if err != nil {
		return nil, err
	}
	if inState != nil {
		err := s.in.LoadState(bytes.NewReader(inState))
		if err != nil {
			return nil, err
		}
	}
	if outState != nil {
		err := s.out.LoadState(bytes.NewReader(outState))
		if err != nil {
			return nil, err
		}
	}
	type container struct {
		pipeCtx *plumber.PipeCtx
		pipe    plumber.Pipe
	}
	if len(s.nodes) < 1 {
		s.nodes = append(s.nodes, pipe.MakePipe(pass))
	}
	var lcs []*container
	for idx, n := range s.nodes {
		lc := &container{
			pipeCtx: &plumber.PipeCtx{Logger: s.logger},
		}
		lc.pipe = n
		if idx == 0 {
			lc.pipeCtx.In, err = s.in.Input()
			if err != nil {
				return nil, err
			}
		} else {
			if lcs == nil {
				return nil, errors.New("error nil lcs array")
			}
			lc.pipeCtx.In = lcs[idx-1].pipeCtx.Out
		}
		if idx == len(s.nodes)-1 {
			lc.pipeCtx.Out, err = s.GetOutputStream().Output()
			if err != nil {
				return nil, err
			}
		} else {
			lc.pipeCtx.Out = make(chan interface{})
		}
		lcs = append(lcs, lc)
	}

	errs := make(chan error, 1024) //TODO: configure error chan cap

	//start checkpoint process
	if s.checkpoint != nil {
		go s.Checkpoint()
	}
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
	s = &defaultPipeline{}
	for _, c := range confs {
		s = c(s)
	}
	return s
}
