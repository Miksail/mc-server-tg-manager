package model

type ServerStatus string

const (
	ServerStatusRunning  ServerStatus = "running"
	ServerStatusStopped  ServerStatus = "stopped"
	ServerStatusStarting ServerStatus = "starting"
	ServerStatusUnknown  ServerStatus = "unknown"
)
