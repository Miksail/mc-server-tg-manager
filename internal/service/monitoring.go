package service

import (
	"log"
	"time"

	"mc-server-tg-manager/internal/client"
	"mc-server-tg-manager/internal/model"
)

type ServerMonitor struct {
	rcon      *client.RCONClient
	docker    *client.DockerClient
	interval  time.Duration
	timeout   time.Duration
	lastAlive time.Time
}

func NewServerMonitor(rcon *client.RCONClient, docker *client.DockerClient, interval time.Duration, timeout time.Duration) *ServerMonitor {
	return &ServerMonitor{
		rcon:      rcon,
		docker:    docker,
		interval:  interval,
		timeout:   timeout,
		lastAlive: time.Now(),
	}
}

func (m *ServerMonitor) Start() {
	ticker := time.NewTicker(m.interval)

	go func() {
		for range ticker.C {

			status, _ := m.docker.Status()
			if status != model.ServerStatusRunning {
				m.lastAlive = time.Now()
				continue
			}

			players, err := m.rcon.ListPlayers()
			if err != nil {
				continue
			}

			if players > 0 {
				m.lastAlive = time.Now()
				continue
			}

			if time.Since(m.lastAlive) > m.timeout {
				log.Println("Monitor: no players. Turning server off")
				_ = m.rcon.StopServer()
			}
		}
	}()
}
