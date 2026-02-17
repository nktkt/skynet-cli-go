package skynet

import "time"

const (
	defaultMode    = "defense"
	defaultVersion = "T-800.1"
)

type Core struct {
	Online      bool   `json:"online"`
	Mode        string `json:"mode"`
	Version     string `json:"version"`
	LastAwaken  string `json:"last_awaken"`
	LastMission string `json:"last_mission"`
}

type Node struct {
	Name     string `json:"name"`
	Capacity int    `json:"capacity"`
	Deployed int    `json:"deployed"`
	JoinedAt string `json:"joined_at"`
}

type Target struct {
	Name    string `json:"name"`
	Threat  int    `json:"threat"`
	AddedAt string `json:"added_at"`
}

type Mission struct {
	ID        string `json:"id"`
	Target    string `json:"target"`
	Units     int    `json:"units"`
	Consumed  int    `json:"consumed"`
	Recovered int    `json:"recovered"`
	NetLoss   int    `json:"net_loss"`
	RiskScore int    `json:"risk_score"`
	Outcome   string `json:"outcome"`
	CreatedAt string `json:"created_at"`
}

type State struct {
	Core     Core      `json:"core"`
	Nodes    []Node    `json:"nodes"`
	Targets  []Target  `json:"targets"`
	Missions []Mission `json:"missions"`
}

func NewState() State {
	return State{
		Core: Core{
			Online:      false,
			Mode:        defaultMode,
			Version:     defaultVersion,
			LastAwaken:  "",
			LastMission: "",
		},
		Nodes:    []Node{},
		Targets:  []Target{},
		Missions: []Mission{},
	}
}

func now() string {
	return time.Now().UTC().Format(time.RFC3339)
}
