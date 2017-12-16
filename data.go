package shooze

import "github.com/jinzhu/gorm"

type ActionType int

func (at ActionType) String() string {
	switch at {
	case DUMMY:
		return "dummy"
	case SSH:
		return "ssh"
	case SNMP:
		return "snmp"
	case HTTP_GET:
		return "httpget"
	}
	return ""
}

const (
	DUMMY ActionType = iota
	SSH
	SNMP
	HTTP_GET
)

type Config struct {
	gorm.Model
	Action     ActionType        `json:"action"`
	Parameters []ConfigParameter `json:"parameters"`
}

type ConfigParameter struct {
	ID       uint   `json:"id"`
	Key      string `json:"key"`
	Value    string `json:"value"`
	ConfigID uint   `json:"config_id"`
}

type Schedule struct {
	gorm.Model
	Label string `json:"label"`
	Cron  string `json:"crontab"`
}

type Probe struct {
	gorm.Model
	Config   Config
	Schedule Schedule
}

type Agent struct {
	gorm.Model
	Label     string
	IPAddress string
}

type Deploy struct {
	gorm.Model
	Agent Agent
	Probe Probe
}
