package graceful

type Process struct {
	starter  starter
	disabled bool
}

func NewProcess(starter starter) Process {
	return Process{
		starter:  starter,
		disabled: false,
	}
}

func (p Process) Disable(d bool) Process {
	p.disabled = d

	return p
}
