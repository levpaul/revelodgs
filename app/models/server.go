package models

import "time"

type Server struct {
	ServerId   int
	UserId     int
	GameId     int
	InstanceId string `json:"-"`
	LaunchTime time.Time
	ExpiryTime time.Time
	State      string
	AmiId      string `json:"-"`
	Options    string
}

const (
	ServerStateRunning string = "RUNNING"
	ServerStateStopped string = "STOPPED"
)
