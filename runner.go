package plumber

type Runner interface {
	ID() int
	SetID(IDFn) Runner
	From(s Stream) Runner
	To(s Stream) Runner
	SetLogger(Logger) Runner
	System() System
	SetSystem(System)
	PartitionBy(Partitioner) Runner
	Run()
}

//IDFn is a function that returns ID for the current Plumber instance.
type IDFn func() int

//Paritioner is a function that gets an event and returns the ID of the handler.
//Partitoner is necessary for clustering Plumber into multiple instances.
//Partitoner should map event key into the ID of plumber instance.
type Partitioner func(event interface{}) int

type runner struct {
	sys         System
	logger      Logger
	partitioner Partitioner
	id          IDFn
}

func (r *runner) PartitionBy(p Partitioner) Runner {
	r.partitioner = p
	return r
}
func (r *runner) SetID(fn IDFn) Runner {
	r.id = fn
	return r
}
func (r *runner) ID() int {
	return r.id()
}
func (r *runner) From(s Stream) Runner {
	r.sys.From(s)
	return r
}

func (r *runner) To(s Stream) Runner {
	r.sys.To(s)
	return r
}

func (r *runner) SetLogger(l Logger) Runner {
	r.logger = l
	return r
}

func (r *runner) System() System {
	return r.sys
}

func (r *runner) Run() {}

func (r *runner) SetSystem(s System) {
	r.sys = s
}

func NewRunner() Runner {
	return &runner{}
}
