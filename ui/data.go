package ui

type programStatusCode int

const (
	statusIdle programStatusCode = iota
	statusConnected
	statusSyncing
)
