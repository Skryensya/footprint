package telemetry

type Status int

const (
	StatusPending Status = iota
	StatusExported
	StatusOrphaned
	StatusSkipped
)
