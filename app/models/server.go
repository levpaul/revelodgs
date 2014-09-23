package models

import "time"

type Server struct {
	UserId     int
	GameId     int
	InstanceId string `json:"-"`
	LaunchTime time.Time
	AmiId      string `json:"-"`
	Options    string
}
