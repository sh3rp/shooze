package shooze

type ActionType int

func (at ActionType) String() string {
	switch at {
	case DUMMY:
		return "dummy"
	case SSH:
		return "ssh"
	case TCP_PING:
		return "tcpping"
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
	TCP_PING
	SNMP
	HTTP_GET
)

type Config struct {
	ID         uint              `json:"id",gorm:"primary_key"`
	Action     ActionType        `json:"action"`
	Parameters []ConfigParameter `json:"parameters"`
	ProbeID    uint              `json:"-"`
}

type ConfigParameter struct {
	ID       uint   `json:"id",gorm:"primary_key"`
	Key      string `json:"key"`
	Value    string `json:"value"`
	ConfigID uint   `json:"-"`
}

type Schedule struct {
	ID      uint   `json:"id",gorm:"primary_key"`
	Label   string `json:"label"`
	Cron    string `json:"crontab"`
	ProbeID uint   `json:"-"`
}

type Probe struct {
	ID       uint     `json:"id",gorm:"primary_key"`
	Config   Config   `json:"config"`
	Schedule Schedule `json:"schedule"`
	DeployID uint     `json:"-"`
}

type Agent struct {
	ID        uint   `json:"id",gorm:"primary_key"`
	Label     string `json:"label"`
	IPAddress string `json:"ip_address"`
	Status    int    `json:"status"`
	DeployID  uint   `json:"-"`
}

type Deploy struct {
	ID     uint  `json:"id",gorm:"primary_key"`
	Agent  Agent `json:"agent"`
	Probe  Probe `json:"probe"`
	Status int   `json:"status"`
}
