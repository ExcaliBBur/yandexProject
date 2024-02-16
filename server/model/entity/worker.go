package entity

import "time"

type Worker struct {
	Hostname      string    `json:"hostname"`
	Id            int       `json:"id"`
	IsAlive       bool      `json:"is_alive"`
	LastHeartBeat time.Time `json:"last_heartbeat"`
}
