package workering

type WorkerStatus int

const (
	_ WorkerStatus = iota
	Running
	Stopped
)
